package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"mksea/common"
	"mksea/output"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/klauspost/compress/zstd"
)

type CompressLevel uint

const (
	CompressHigh CompressLevel = iota
	CompressMid
	CompressLow
)

func (l CompressLevel) ToZstdLevel() zstd.EncoderLevel {
	switch l {
	case CompressHigh:
		return zstd.SpeedBestCompression
	case CompressMid:
		return zstd.SpeedBetterCompression
	case CompressLow:
		return zstd.SpeedDefault
	default:
		panic("unreachable")
	}
}

func (l *CompressLevel) FromString(s string) bool {
	switch s {
	case "high":
		*l = CompressHigh
	case "mid":
		*l = CompressMid
	case "low":
		*l = CompressLow
	default:
		return false
	}
	return true
}

type Packer struct {
	CompressLevel CompressLevel
	Output        string
	BaseName      string
	Input         []string
	Encrypt       bool
	Platforms     []TargetPlatform
	Threads       int
	Gui           bool
	Silent        bool

	archiveName string
	archiveSize int64
}

func (p *Packer) Pack() error {
	if len(p.Platforms) == 0 {
		var tp TargetPlatform
		tp.FromString("")
		p.Platforms = append(p.Platforms, tp)
	}
	if err := p.archive(); err != nil {
		return common.NewContextError("archive failed", err)
	}
	if err := p.generate(); err != nil {
		return common.NewContextError("self extractable archive generating failed", err)
	}
	if err := p.build(); err != nil {
		return common.NewContextError("self extractable archive building failed", err)
	}
	p.log("> done!")
	return nil
}

func (p *Packer) archive() error {
	if p.Encrypt {
		if err := p.generateEncryptionKey(); err != nil {
			return common.NewContextError("generating encryption key", err)
		}
	}

	output.Env.EncoderLevel = p.CompressLevel.ToZstdLevel()
	if p.Threads == 0 {
		output.Env.EncoderThreads = runtime.NumCPU()
	} else {
		output.Env.EncoderThreads = p.Threads
	}

	archiveName := "archive.dat"
	p.archiveName = filepath.Join(workInstallerDir, archiveName)
	archiveName, _ = filepath.Rel(workDir, p.archiveName)
	out, err := output.Open(archiveName, 0644)
	if err != nil {
		return common.NewContextError("cannot open archive for write", err)
	}
	closeOutput := func() error {
		if err := out.Close(); err != nil {
			return common.NewContextError("close output error", err)
		}
		return nil
	}
	errList := common.NewErrorListCap(3)

	p.log("> packing files...")
	l := len(p.Input)
	for i, it := range p.Input {
		rel, _ := filepath.Rel(workDir, it)
		p.logf("> [%d/%d] packing \"%s\"...", i+1, l, rel)
		outFile, err := out.Next(it)
		if err != nil {
			errList.Append(
				common.NewContextError("cannot open output file write", err),
				closeOutput(),
			)
			return errList
		}
		closeOutputFile := func() error {
			if err := outFile.Close(); err != nil {
				return common.NewContextError("close output file error", err)
			}
			return nil
		}
		inFile, err := os.Open(it)
		if err != nil {
			errList.Append(
				common.NewContextError("cannot open input file for read", err),
				closeOutputFile(),
				closeOutput(),
			)
			return errList
		}
		closeInputFile := func() error {
			if err := inFile.Close(); err != nil {
				return common.NewContextError("close input file error", err)
			}
			return nil
		}
		_, err = io.Copy(outFile, inFile)
		if err != nil {
			errList.Append(
				common.NewContextError("cannot pack file", err),
			)
		}
		errList.Append(
			closeInputFile(),
			closeOutputFile(),
		)
		if err != nil {
			errList.Append(
				closeOutput(),
			)
			return errList.RealError()
		}
	}

	if err := closeOutput(); err != nil {
		return err
	}

	info, err := os.Stat(p.archiveName)
	if err != nil {
		return common.NewContextError("cannot obtain archive size", err)
	}
	p.archiveSize = info.Size()
	return nil
}

func (p *Packer) generateEncryptionKey() error {
	key := make([]byte, 128)
	for i := range key {
		n := rand.Uint64()
		for j := 1; j < 8; j++ {
			n ^= (n >> (8 * j) & 0xff)
		}
		key[i] = byte(n & 0xff)
	}
	output.Env.Encode = true
	output.Env.EncoderKey = key
	return nil
}

func (p *Packer) generate() error {
	writeBytes := func(name string, data []byte) error {
		target := filepath.Join(workInstallerDir, name)
		if err := os.MkdirAll(filepath.Dir(target), 0700); err != nil {
			return err
		}
		return os.WriteFile(target, data, 0600)
	}
	write := func(name, data string) error {
		return writeBytes(name, []byte(data))
	}

	p.log("> generating self extractable archive...")

	if err := unpackTemplate(); err != nil {
		return common.NewContextError("template unpacking", err)
	}

	if err := write("cmd/sea/archive_init.go", fmt.Sprintf(`package main

func init() {
	archiveSize = %d
}
`, p.archiveSize)); err != nil {
		return common.NewContextError("archive info generating", err)
	}

	if err := write("cmd/sea/version_init.go", fmt.Sprintf(`package main

func init() {
	Version.Major = %d
	Version.Minor = %d
	Version.Patch = %d
	Version.Suffix = "%s"
}
`, Version.Major, Version.Minor, Version.Patch, Version.Suffix)); err != nil {
		return common.NewContextError("version generating", err)
	}

	if p.Encrypt {
		if err := writeBytes("cmd/sea/encoder.key", output.Env.EncoderKey); err != nil {
			return err
		}

		if err := write("cmd/sea/decoder_init.go", `package main

import (
	_ "embed"
	"mksea/input"
)

//go:embed encoder.key
var decodeKey []byte

func init() {
	input.Env.Decode = true
	input.Env.DecodeKey = decodeKey
}
`); err != nil {
			return common.NewContextError("decoder generating", err)
		}
	}

	if err := p.goMod(); err != nil {
		return common.NewContextError("module initialization", err)
	}

	return nil
}

func (p *Packer) build() error {
	var buildFunc func(string, string, TargetPlatform) error
	if p.Gui {
		if _, err := exec.LookPath("fyne-cross"); err != nil {
			return errors.New("fyne-cross not found")
		}
		buildFunc = p.buildFyne
	} else {
		buildFunc = p.buildCli
	}
	l := len(p.Platforms)
	p.log("> building executables...")
	pkg := filepath.ToSlash(filepath.Join(workInstallerDir, "cmd/sea"))
	if len(p.BaseName) == 0 {
		p.BaseName = filepath.Base(workDir)
	}
	for i, it := range p.Platforms {
		baseName := fmt.Sprintf("%s_%s_%s", p.BaseName, it.OsName(), it.ArchName())
		if it.OsName() == "windows" {
			baseName += ".exe"
		}
		p.logf("> [%d/%d] build %s...", i+1, l, baseName)
		target := filepath.Join(workInstallerDir, baseName)
		if err := buildFunc(pkg, target, it); err != nil {
			return err
		}
		if err := p.join(baseName, target); err != nil {
			return common.NewContextError("join failed", err)
		}
	}
	return nil
}

func (p *Packer) buildCli(pkg, target string, platform TargetPlatform) error {
	cmd := exec.Command(
		"go",
		"build",
		"-v",
		"-o", target,
		pkg,
	)
	cmd.Dir = workInstallerDir
	cmd.Env = append(cmd.Env, "GOOS="+platform.OsName(), "GOARCH="+platform.ArchName())
	return runCommand(cmd)
}

func (p *Packer) buildFyne(pkg, target string, platform TargetPlatform) error {
	targetName := filepath.Base(target)
	// TODO: icon
	cmd := exec.Command(
		"fyne-cross",
		platform.OsName(),
		"-arch", platform.ArchName(),
		"-name", targetName,
		"-app-id", "com.github.egor9814.mksea",
		"-tags", "fyne_gui",
		"-no-cache",
		pkg,
	)
	cmd.Env = os.Environ()
	cmd.Dir = workInstallerDir
	return runCommand(cmd)
}

func (p *Packer) goMod() error {
	errList := common.NewErrorList()
	errList.Append(runGoModCommand("init", "mksea"))
	if errList.Len() == 0 {
		errList.Append(runGoModCommand("tidy"))
	}
	return errList.RealError()
}

func runGoModCommand(args ...string) error {
	newArgs := make([]string, 1+len(args))
	newArgs[0] = "mod"
	for i, it := range args {
		newArgs[i+1] = it
	}
	cmd := exec.Command("go", newArgs...)
	cmd.Dir = workInstallerDir
	return runCommand(cmd)
}

func runCommand(cmd *exec.Cmd) error {
	cmd.Env = append(cmd.Env, "GOPATH="+goPath, "GOCACHE="+goCache, "GOTMPDIR="+goTmp)
	var out bytes.Buffer
	cmd.Stderr = &out
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		err = errors.New(out.String())
	}
	return err
}

func (p *Packer) join(outputName, inputName string) error {
	outFile, err := output.OpenRaw(filepath.Join(p.Output, outputName), 0755)
	if err != nil {
		return err
	}
	errList := common.NewErrorList()
	_, err = copyFile(outFile, inputName)
	errList.Append(err)
	if err == nil {
		_, err = copyFile(outFile, p.archiveName)
		errList.Append(err)
	}
	errList.Append(outFile.Close())
	return errList.RealError()
}

func copyFile(o io.Writer, inputName string) (int64, error) {
	inputFile, err := os.Open(inputName)
	if err != nil {
		return 0, err
	}
	n, copyErr := io.Copy(o, inputFile)
	return n, common.NewErrorListFrom(copyErr, inputFile.Close()).RealError()
}

func (p *Packer) log(args ...any) {
	if !p.Silent {
		log.Println(args...)
	}
}

func (p *Packer) logf(format string, args ...any) {
	if !p.Silent {
		log.Printf(format, args...)
	}
}

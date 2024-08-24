//go:build !tui && !gui && !dui

package main

import (
	"fmt"
	"io"
	"log"
	"mksea/input"
	"mksea/output"
	"os"
	"path/filepath"
)

func run(args []string) {
	targetPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("cannot obtain work dir: %v", err)
	}
	input.Env.MaxMem = 0
	for i, l := 1, len(args); i < l; i++ {
		arg := args[i]
		switch arg {
		case "--show-exe-size":
			_, exeSize, err := archiveOffset()
			if err != nil {
				log.Fatal(err)
			}
			log.Println(exeSize)
			return

		case "--show-arc-size":
			log.Println(archiveSize)
			return

		case "-h", "--help":
			printHelp()
			return

		case "-v", "--version":
			printVersion()
			return

		case "-o", "--output":
			i++
			if i >= l {
				log.Fatalf("output dir not specified after %q flag", arg)
			}
			newTargetPath := args[i]
			startPoint := newTargetPath
			for {
				if info, err := os.Stat(startPoint); err != nil {
					startPoint = filepath.Dir(startPoint)
				} else if !info.IsDir() {
					log.Fatalf("path \"%s\" is not a directory", startPoint)
				} else {
					break
				}
			}
			if err := os.MkdirAll(newTargetPath, 0755); err != nil {
				log.Fatalf("cannot create output directory: %v", err)
			}
			targetPath = newTargetPath

		case "-m", "--max":
			i++
			if i >= l {
				log.Fatalf("max memory limit not specified after %q flag", arg)
			}
			newMemoryLimit := uint64(0)
			r := []rune(args[i])
			i, l := 0, len(r)
			for ; i < l; i++ {
				if '0' <= r[i] && r[i] <= '9' {
					newMemoryLimit = newMemoryLimit*10 + uint64(r[i]-'0')
				} else {
					break
				}
			}
			if i == 0 {
				log.Fatalf("invalid memory limit value %q", args[i])
			}
			if i < l {
				valid := false
				for _, it := range "KMGT" {
					newMemoryLimit *= 1024
					if it == r[i] {
						valid = true
						break
					}
				}
				if !valid {
					log.Fatalf("unsupported memory limit suffix %q", r[i:])
				}
			}
			input.Env.MaxMem = newMemoryLimit

		default:
			log.Fatalf("unsupported argument %q, type '%s --help' for get help", arg, args[0])
		}
	}

	if input.Env.MaxMem == 0 {
		input.Env.MaxMem = 8 * 1024 * 1024 * 1024 // 8GB by default
	}
	output.Env.WorkDir = targetPath

	exe, exeOffset, err := archiveOffset()
	if err != nil {
		log.Fatalf("cannot obtain info from sea file: %v", err)
	}
	in, err := input.Open(exe, exeOffset)
	if err != nil {
		log.Fatalf("cannot open sea file: %v", err)
	}
	defer in.Close()
	progress := func() int {
		p := in.Progress()
		return int(100.0 * float64(p.Current()) / float64(p.All()))
	}

	log.Println("[  0%] preparing...")
	for it, err := in.Next(); it != nil || err != nil; it, err = in.Next() {
		if err != nil {
			log.Fatalf("cannot go to next file: %v", err)
		}
		log.Printf(
			"[%3d%%] unpacking \"%s\"...\n",
			progress(),
			filepath.FromSlash(it.Path),
		)
		outFile, err := output.OpenRaw(it.Path, 0755)
		if err != nil {
			log.Fatalf("cannot open file for write: %v", err)
		}

		if _, err := io.Copy(outFile, it.Reader); err != nil {
			outFile.Close()
			log.Fatalf("cannot unpack file: %v", err)
		}

		outFile.Close()
	}
	log.Println("[100%] done!")
}

func main() {
	run(os.Args)
}

func printHelp() {
	fmt.Printf("usage of Self Extractable Archive: %s [options...]\n", os.Args[0])
	fmt.Printf("options:\n")
	fmt.Printf("  -h, --help      - print this help\n")
	fmt.Printf("  -v, --version   - print sea version\n")
	fmt.Printf("  -o, --output    - set output directory (current dir by default)\n")
	fmt.Printf("  -m, --max       - set memory limit (8G by default)\n")
	fmt.Println()
}

func printVersion() {
	fmt.Printf("Self Extractable Archive v%d.%d.%d%s", Version.Major, Version.Minor, Version.Patch, Version.Suffix)
	fmt.Println()
}

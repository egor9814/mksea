package main

import (
	"io"
	"log"
	"mksea/output"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/zstd"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("cannot obtain work dir: %v", err)
	}
	name := filepath.Join(wd, "template.tar.zst")
	wd = filepath.Dir(filepath.Dir(wd))

	files := []string{
		"cmd/sea/archive.go",
		"cmd/sea/main_gui.go",
		"cmd/sea/main_no_gui.go",
		"cmd/sea/main.go",
		"cmd/sea/meta.go",
		"cmd/sea/password.go",
		"cmd/sea/version.go",
		"common/errors.go",
		"common/meta.go",
		"common/password.go",
		"crypto/xor.go",
		"input/decoder.go",
		"input/env.go",
		"input/interface.go",
		"input/open.go",
		"input/raw_input.go",
		"input/tar_input.go",
		"input/zstd.go",
		"output/env.go",
		"output/open_raw.go",
		"output/raw_output.go",
	}

	output.Env.EncoderThreads = 1
	output.Env.EncoderLevel = zstd.SpeedBestCompression

	out, err := output.Open(name, 0644)
	if err != nil {
		log.Fatalf("cannot open output file: %v", err)
	}
	defer out.Close()

	output.Env.WorkDir = wd
	for _, it := range files {
		p := filepath.Join(wd, it)
		outFile, err := out.Next(p)
		if err != nil {
			log.Fatal(err)
		}
		inFile, err := os.Open(p)
		if err != nil {
			log.Fatal(err)
		}
		io.Copy(outFile, inFile)
		inFile.Close()
		outFile.Close()
	}
}

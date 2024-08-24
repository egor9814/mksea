package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var workDir string

var workInstallerDir string

var goPath, goCache, goTmp string

func init() {
	if wd, err := os.Getwd(); err != nil {
		log.Fatalf("cannot obtain work dir: %v", err)
	} else {
		workDir = wd
	}
	workInstallerDir = filepath.Join(workDir, "mksea-temp")

	goCache = filepath.Join(workInstallerDir, ".cache")
	goTmp = filepath.Join(workInstallerDir, ".tmp")
	if err := os.MkdirAll(goTmp, 0755); err != nil {
		log.Fatalf("cannot create go temp directory: %v", err)
	}

	found := false
	goPath, found = os.LookupEnv("GOPATH")
	if !found {
		cmd := exec.Command("go", "env", "GOPATH")
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err == nil {
			goPath = strings.TrimSpace(out.String())
			info, err := os.Stat(goPath)
			if err == nil && info.IsDir() {
				found = true
			}
		}
	}
	if !found {
		log.Fatal("GOPATH not provided")
	}
}

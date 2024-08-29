package main

import (
	"bytes"
	"errors"
	"mksea/common"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var workDir string

var workInstallerDir string

var goPath, goCache, goTmp string

func init_wd() error {
	if wd, err := os.Getwd(); err != nil {
		return common.NewContextError("cannot obtain work dir", err)
	} else {
		workDir = wd
	}
	workInstallerDir = filepath.Join(workDir, "mksea-temp")

	goCache = filepath.Join(workInstallerDir, ".cache")
	goTmp = filepath.Join(workInstallerDir, ".tmp")
	if err := os.MkdirAll(goTmp, 0755); err != nil {
		return common.NewContextError("cannot create go temp directory", err)
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
		return errors.New("GOPATH not provided")
	}

	return nil
}

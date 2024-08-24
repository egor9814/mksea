package main

import (
	"mksea/common"
	"os"
	"path/filepath"
)

var archiveSize int64
var archiveIsTest bool

func archiveOffset() (string, int64, error) {
	if archiveIsTest {
		return "archive.dat", 0, nil
	}
	exe, err := os.Executable()
	if err != nil {
		return "", 0, common.NewContextError("cannot retrieve sea path", err)
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return "", 0, common.NewContextError("cannot resolve symlinks for sea path", err)
	}
	info, err := os.Stat(exe)
	if err != nil {
		return "", 0, common.NewContextError("cannot obtain executable info", err)
	}
	return exe, info.Size() - archiveSize, nil
}

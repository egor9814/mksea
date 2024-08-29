package main

import (
	"mksea/common"
	"os"
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
	info, err := os.Stat(exe)
	if err != nil {
		return "", 0, common.NewContextError("cannot obtain executable info", err)
	}
	return exe, info.Size() - archiveSize, nil
}

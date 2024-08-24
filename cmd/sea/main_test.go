package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	archiveIsTest = true
	defer func() {
		if reason := recover(); reason != nil {
			t.Fatal(reason)
		}
	}()
	run([]string{
		os.Args[0],
		"-o", "tempdir",
	})
}

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
	new_cli_app().Run([]string{
		os.Args[0],
		"-o", "tempdir",
	})
}

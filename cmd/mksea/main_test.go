package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	if err := new_cli_app().Run([]string{
		os.Args[0],
		// "-e",
		// "-p", "windows/amd64",
		// "-p", "linux",
		// "*.go",
		// "testdata/b*/*a*",
		// "testdata/*",
		"testdata",
	}); err != nil {
		t.Fatal(err)
	}
}

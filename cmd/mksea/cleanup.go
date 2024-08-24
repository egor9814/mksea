package main

import "os"

func cleanup() error {
	return os.RemoveAll(workInstallerDir)
}

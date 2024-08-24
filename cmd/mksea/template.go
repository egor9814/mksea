package main

import (
	_ "embed"
	"io"
	"mksea/input"
	"mksea/output"
)

//go:generate go run ../template/.

//go:embed template.tar.zst
var tempalteArchive []byte

func unpackTemplate() error {
	in, err := input.OpenBytes(tempalteArchive)
	if err != nil {
		return err
	}
	defer in.Close()

	oldWd := output.Env.WorkDir
	output.Env.WorkDir = workInstallerDir
	defer func() {
		output.Env.WorkDir = oldWd
	}()

	for it, err := in.Next(); it != nil || err != nil; it, err = in.Next() {
		if err != nil {
			return err
		}
		out, err := output.OpenRaw(it.Path, 0644)
		if err != nil {
			return nil
		}
		if _, err := io.Copy(out, it.Reader); err != nil {
			return err
		}
		out.Close()
	}

	return nil
}

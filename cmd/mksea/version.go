package main

var Version struct {
	Major, Minor, Patch int
	Suffix              string
}

//go:generate go run ../version/.

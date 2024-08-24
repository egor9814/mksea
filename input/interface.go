package input

import "io"

type File struct {
	Path   string
	Reader io.Reader
}

type ProgressStatus interface {
	Current() int64
	All() int64
}

type Interface interface {
	io.Closer
	Next() (*File, error)
	Progress() ProgressStatus
}

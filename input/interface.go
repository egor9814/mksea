package input

import "io"

type ArchiveFormat uint

type File struct {
	Path   string
	Reader io.Reader
}

type Interface interface {
	io.Closer
	Next() (*File, error)
}

package output

import "io"

type Interface interface {
	io.Closer
	Next(name string) (io.WriteCloser, error)
}

package input

import (
	"bytes"
	"os"
)

type bytesInput struct {
	r    *bytes.Reader
	size int64
	read int64
}

func (i *bytesInput) Close() error {
	i.r = nil
	return nil
}

func (i *bytesInput) Read(b []byte) (n int, err error) {
	if i.r == nil {
		return 0, os.ErrClosed
	}
	n, err = i.r.Read(b)
	i.read += int64(n)
	return
}

func newBytesInput(data []byte) *bytesInput {
	return &bytesInput{
		r:    bytes.NewReader(data),
		size: int64(len(data)),
	}
}

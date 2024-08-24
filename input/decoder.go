package input

import (
	"io"
	"os"
)

type decoderInput struct {
	reader io.Reader
	closer io.Closer
	index  int
}

func (i *decoderInput) Read(b []byte) (n int, err error) {
	if i.reader == nil {
		return 0, os.ErrClosed
	}
	n, err = i.reader.Read(b)
	for j, it := range b[:n] {
		b[j] = it ^ Env.DecodeKey[i.index]
		i.index = (i.index + 1) % len(Env.DecodeKey)
	}
	return
}

func (i *decoderInput) Close() (err error) {
	if i.closer != nil {
		err = i.closer.Close()
		i.closer = nil
	}
	i.reader = nil
	return
}

func newDecoderInput(r io.Reader, c io.Closer) *decoderInput {
	return &decoderInput{
		reader: r,
		closer: c,
	}
}

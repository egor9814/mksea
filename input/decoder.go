package input

import (
	"io"
	"mksea/crypto"
	"os"
)

type decoderInput struct {
	reader *crypto.XorReader
	closer io.Closer
}

func (i *decoderInput) Read(b []byte) (n int, err error) {
	if i.reader == nil {
		return 0, os.ErrClosed
	}
	n, err = i.reader.Read(b)
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
		reader: crypto.NewXorReader(r, Env.DecodeKey),
		closer: c,
	}
}

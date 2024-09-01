package input

import (
	"io"

	"github.com/klauspost/compress/zstd"
)

func newZstdInput(r io.Reader, c io.Closer) (Interface, error) {
	z, err := zstd.NewReader(
		r,
		zstd.WithDecoderConcurrency(0),
		zstd.WithDecoderLowmem(false),
		zstd.WithDecoderMaxMemory(Env.MaxMem),
	)
	if err != nil {
		return nil, err
	}
	t := newTarInput(z, c)
	return t, nil
}

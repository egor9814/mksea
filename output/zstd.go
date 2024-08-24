package output

import (
	"fmt"
	"io"

	"github.com/klauspost/compress/zstd"
)

type zstdCloser struct {
	e *zstd.Encoder
	c io.Closer
}

func (c *zstdCloser) Close() (err error) {
	if c.e != nil {
		err = c.e.Close()
		err2 := c.c.Close()
		if err2 != nil {
			if err != nil {
				err = fmt.Errorf("zstd output close error:\n  %v\n  %v", err, err2)
			} else {
				err = err2
			}
		}
		c.e = nil
		c.c = nil
	}
	return
}

func newZstdOutput(w io.Writer, c io.Closer) (Interface, error) {
	z, err := zstd.NewWriter(
		w,
		zstd.WithEncoderLevel(Env.EncoderLevel),
		zstd.WithEncoderConcurrency(Env.EncoderThreads),
	)
	if err != nil {
		return nil, err
	}
	return newTarOutput(z, &zstdCloser{
		e: z,
		c: c,
	}), nil
}

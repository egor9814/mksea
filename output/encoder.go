package output

import (
	"io"
	"os"
)

type encoderOutput struct {
	writer io.Writer
	closer io.Closer
	index  int
}

func (o *encoderOutput) Write(b []byte) (int, error) {
	if o.writer == nil {
		return 0, os.ErrClosed
	}
	b2 := make([]byte, len(b))
	for i, it := range b {
		b2[i] = it ^ Env.EncoderKey[o.index]
		o.index = (o.index + 1) % len(Env.EncoderKey)
	}
	return o.writer.Write(b2)
}

func (o *encoderOutput) Close() (err error) {
	if o.closer != nil {
		err = o.closer.Close()
		o.closer = nil
	}
	o.writer = nil
	return
}

func newEncoderOutput(w io.Writer, c io.Closer) *encoderOutput {
	return &encoderOutput{
		writer: w,
		closer: c,
	}
}

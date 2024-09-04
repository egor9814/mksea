package output

import (
	"io"
	"mksea/crypto"
	"os"
)

type encoderOutput struct {
	writer *crypto.XorWriter
	closer io.Closer
}

func (o *encoderOutput) Write(b []byte) (int, error) {
	if o.writer == nil {
		return 0, os.ErrClosed
	}
	return o.writer.Write(b)
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
		writer: crypto.NewXorWriter(w, Env.EncoderKey),
		closer: c,
	}
}

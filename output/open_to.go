package output

import (
	"io"

	"github.com/klauspost/compress/zstd"
)

func OpenTo(wc io.WriteCloser) (Interface, error) {
	return OpenTo2(wc, wc)
}

func OpenTo2(w io.Writer, c io.Closer) (Interface, error) {
	if Env.Encode {
		wc := newEncoderOutput(w, c)
		w = wc
		c = wc
	}
	if Env.EncoderLevel == (zstd.EncoderLevel)(0) {
		return newTarOutput(w, c), nil
	}
	return newZstdOutput(w, c)
}

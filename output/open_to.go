package output

import (
	"io"
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
	if Env.XzEncode {
		return newXzOutput(w, c)
	}
	if Env.ZstdEncoderLevel != ZstdEncoderLevelNone {
		return newZstdOutput(w, c)
	}
	return newTarOutput(w, c), nil
}

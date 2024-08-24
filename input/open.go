package input

import (
	"io"
)

func Open(name string, offset int64) (Interface, error) {
	raw, err := newRawInput(name, offset)
	if err != nil {
		return nil, err
	}
	var rc io.ReadCloser = raw
	if Env.Decode {
		rc = newDecoderInput(rc, rc)
	}
	return newZstdInput(rc, rc, raw.Progress())
}

package input

import "io"

func OpenBytes(data []byte) (Interface, error) {
	b := newBytesInput(data)
	var rc io.ReadCloser = b
	if Env.Decode {
		rc = newDecoderInput(rc, rc)
	}
	return newZstdInput(rc, rc, b.Progress())
}

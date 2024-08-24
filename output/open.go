package output

import (
	"io"
	"os"
)

func Open(name string, mode os.FileMode) (Interface, error) {
	raw, err := OpenRaw(name, mode)
	if err != nil {
		return nil, err
	}
	var wc io.WriteCloser = raw
	if Env.Encode {
		wc = newEncoderOutput(wc, wc)
	}
	return newZstdOutput(wc, wc)
}

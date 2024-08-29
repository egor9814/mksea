package output

import (
	"io"
	"os"

	"github.com/klauspost/compress/zstd"
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
	if Env.EncoderLevel == (zstd.EncoderLevel)(0) {
		return newTarOutput(wc, wc), nil
	}
	return newZstdOutput(wc, wc)
}

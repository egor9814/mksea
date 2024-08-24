package output

import (
	"io"
	"os"
)

func OpenRaw(name string, mode os.FileMode) (io.WriteCloser, error) {
	return newRawOutput(name, mode)
}

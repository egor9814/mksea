package output

import (
	"os"
)

func Open(name string, mode os.FileMode) (Interface, error) {
	raw, err := OpenRaw(name, mode)
	if err != nil {
		return nil, err
	}
	return OpenTo(raw)
}

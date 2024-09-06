package input

import (
	"io"

	"github.com/ulikunitz/xz"
)

func newXzInput(r io.Reader, c io.Closer) (Interface, error) {
	xzr, err := xz.NewReader(r)
	if err != nil {
		return nil, err
	}
	return newTarInput(xzr, c), nil
}

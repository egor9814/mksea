package output

import (
	"io"
	"mksea/common"

	"github.com/ulikunitz/xz"
)

type xzCloser struct {
	w *xz.Writer
	c io.Closer
}

func (c *xzCloser) Close() error {
	if c.w == nil && c.c == nil {
		return nil
	}
	errList := common.NewErrorList()
	if c.w != nil {
		errList.Append(c.w.Close())
		c.w = nil
	}
	if c.c != nil {
		errList.Append(c.c.Close())
		c.c = nil
	}
	return errList.RealError()
}

func newXzOutput(w io.Writer, c io.Closer) (Interface, error) {
	xzw, err := xz.NewWriter(w)
	if err != nil {
		return nil, err
	}
	return newTarOutput(xzw, &xzCloser{
		w: xzw,
		c: c,
	}), nil
}

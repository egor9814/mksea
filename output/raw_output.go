package output

import (
	"os"
	"path/filepath"
)

type rawOutput struct {
	f *os.File
}

func (o *rawOutput) Write(b []byte) (int, error) {
	if o.f == nil {
		return 0, os.ErrClosed
	}
	return o.f.Write(b)
}

func (o *rawOutput) Close() (err error) {
	if o.f != nil {
		err = o.f.Close()
		o.f = nil
	}
	return
}

func newRawOutput(name string, mode os.FileMode) (*rawOutput, error) {
	p := filepath.FromSlash(filepath.Join(Env.WorkDir, name))
	d := filepath.Dir(p)
	if err := os.MkdirAll(d, 0755); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return nil, err
	}
	return &rawOutput{f}, nil
}

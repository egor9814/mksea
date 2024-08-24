package output

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type tarOutput struct {
	writer *tar.Writer
	closer io.Closer
}

type tarOutputWriter struct {
	o *tarOutput
}

func (w *tarOutputWriter) Write(b []byte) (int, error) {
	if w.o == nil {
		return 0, os.ErrClosed
	}
	return w.o.Write(b)
}

func (w *tarOutputWriter) Close() error {
	w.o = nil
	return nil
}

func (o *tarOutput) Close() (err error) {
	if o.closer != nil {
		err = o.writer.Close()
		err2 := o.closer.Close()
		o.closer = nil
		o.writer = nil
		if err2 != nil {
			if err != nil {
				err = fmt.Errorf("tar close error:\n  %v\n  %v", err, err2)
			} else {
				err = err2
			}
		}
	}
	return
}

func (o *tarOutput) Write(b []byte) (int, error) {
	if o.writer == nil {
		return 0, os.ErrClosed
	}
	return o.writer.Write(b)
}

func (o *tarOutput) Next(name string) (io.WriteCloser, error) {
	p, err := filepath.Rel(Env.WorkDir, name)
	if err != nil {
		return nil, err
	}
	info, _ := os.Stat(name)
	header := tar.Header{
		Name: filepath.ToSlash(p),
		Mode: int64(info.Mode()),
		Size: info.Size(),
	}
	if err := o.writer.WriteHeader(&header); err != nil {
		return nil, err
	}
	return &tarOutputWriter{o}, nil
}

func newTarOutput(w io.Writer, c io.Closer) *tarOutput {
	return &tarOutput{
		writer: tar.NewWriter(w),
		closer: c,
	}
}

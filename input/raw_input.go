package input

import (
	"fmt"
	"io"
	"os"
)

type rawInput struct {
	f    *os.File
	size int64
	read int64
}

func (i *rawInput) Close() (err error) {
	if i.f != nil {
		err = i.f.Close()
		i.f = nil
	}
	return
}

func (i *rawInput) Read(b []byte) (n int, err error) {
	if i.f == nil {
		err = os.ErrClosed
		return
	}
	n, err = i.f.Read(b)
	i.read += int64(n)
	return
}

func (i *rawInput) open(name string, offset int64) (err error) {
	if err = i.Close(); err != nil {
		return fmt.Errorf("cannot reopen input: %v", err)
	}
	if info, err := os.Stat(name); err != nil {
		return err
	} else {
		i.size = info.Size() - offset
	}
	i.f, err = os.Open(name)
	if err == nil {
		_, err = i.f.Seek(offset, io.SeekStart)
	}
	return
}

func newRawInput(name string, offset int64) (*rawInput, error) {
	i := new(rawInput)
	return i, i.open(name, offset)
}

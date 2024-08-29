package input

import (
	"fmt"
	"io"
	"os"
)

type rawInput struct {
	f        *os.File
	size     int64
	read     int64
	progress chan int64
}

func (i *rawInput) Close() (err error) {
	if i.progress != nil {
		close(i.progress)
		i.progress = nil
	}
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
	if i.progress != nil {
		i.progress <- i.read
	}
	return
}

func (i *rawInput) Progress() ProgressStatus {
	return i
}

func (i *rawInput) Current() int64 {
	return i.read
}

func (i *rawInput) All() int64 {
	return i.size
}

func (i *rawInput) Chan() <-chan int64 {
	if i.progress == nil {
		i.progress = make(chan int64)
	}
	return i.progress
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

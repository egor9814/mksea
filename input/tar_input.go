package input

import (
	"archive/tar"
	"io"
	"os"
)

type tarInput struct {
	reader *tar.Reader
	closer io.Closer
}

type tarFileReader struct {
	i *tarInput
}

func (r *tarFileReader) Read(b []byte) (int, error) {
	if r.i.reader == nil {
		return 0, os.ErrClosed
	}
	return r.i.reader.Read(b)
}

func (r *tarFileReader) Close() error {
	r.i.reader = nil
	return nil
}

func (i *tarInput) Next() (*File, error) {
	if i.reader == nil {
		return nil, os.ErrClosed
	}
	header, err := i.reader.Next()
	if err == io.EOF {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &File{
		Path:   header.Name,
		Reader: &tarFileReader{i},
	}, nil
}

func (i *tarInput) Close() (err error) {
	if i.closer != nil {
		err = i.closer.Close()
		i.closer = nil
	}
	i.reader = nil
	return
}

func newTarInput(r io.Reader, c io.Closer) *tarInput {
	return &tarInput{
		reader: tar.NewReader(r),
		closer: c,
	}
}

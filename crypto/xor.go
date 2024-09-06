package crypto

import "io"

type XorKey []byte

type XorKeyReader struct {
	data  XorKey
	index int
}

func (k *XorKeyReader) IsValid() bool {
	return len(k.data) > 0
}

func (k *XorKeyReader) Next() (b byte) {
	b = k.data[k.index]
	k.index = (k.index + 1) % len(k.data)
	return
}

func (k *XorKeyReader) Reset(data XorKey) *XorKeyReader {
	k.data = data
	k.index = 0
	return k
}

func (k *XorKeyReader) ResetPosition() {
	k.index = 0
}

func NewXorKeyReader(data XorKey) *XorKeyReader {
	return new(XorKeyReader).Reset(data)
}

type XorReader struct {
	r io.Reader
	k XorKeyReader
}

func (r *XorReader) Reset(reader io.Reader, key XorKey) *XorReader {
	r.r = reader
	r.k.Reset(key)
	return r
}

func (r *XorReader) Read(p []byte) (n int, err error) {
	if r.r == nil {
		err = io.ErrClosedPipe
		return
	}
	n, err = r.r.Read(p)
	if r.k.IsValid() && n > 0 {
		for i, it := range p[:n] {
			p[i] = it ^ r.k.Next()
		}
	}
	return
}

func NewXorReader(reader io.Reader, key XorKey) *XorReader {
	return new(XorReader).Reset(reader, key)
}

type XorWriter struct {
	w io.Writer
	k XorKeyReader
}

func (w *XorWriter) Reset(writer io.Writer, key XorKey) *XorWriter {
	w.w = writer
	w.k.Reset(key)
	return w
}

func (w *XorWriter) Write(p []byte) (n int, err error) {
	if w.w == nil {
		err = io.ErrClosedPipe
		return
	}
	if w.k.IsValid() {
		b := make([]byte, len(p))
		for i, it := range p {
			b[i] = it ^ w.k.Next()
		}
		p = b
	}
	n, err = w.w.Write(p)
	return
}

func NewXorWriter(writer io.Writer, key XorKey) *XorWriter {
	return new(XorWriter).Reset(writer, key)
}

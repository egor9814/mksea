package crypto

import "io"

type XorKey []byte

type xorKeyReader struct {
	data  XorKey
	index int
}

func (k *xorKeyReader) isValid() bool {
	return len(k.data) > 0
}

func (k *xorKeyReader) next() (b byte) {
	b = k.data[k.index]
	k.index = (k.index + 1) % len(k.data)
	return
}

func (k *xorKeyReader) reset(data XorKey) {
	k.data = data
	k.index = 0
}

type XorReader struct {
	r io.Reader
	k xorKeyReader
}

func (r *XorReader) Reset(reader io.Reader, key XorKey) *XorReader {
	r.r = reader
	r.k.reset(key)
	return r
}

func (r *XorReader) Read(p []byte) (n int, err error) {
	if r.r == nil {
		err = io.ErrClosedPipe
		return
	}
	n, err = r.r.Read(p)
	if r.k.isValid() && n > 0 {
		for i, it := range p[:n] {
			p[i] = it ^ r.k.next()
		}
	}
	return
}

func NewXorReader(reader io.Reader, key XorKey) *XorReader {
	return new(XorReader).Reset(reader, key)
}

type XorWriter struct {
	w io.Writer
	k xorKeyReader
}

func (w *XorWriter) Reset(writer io.Writer, key XorKey) *XorWriter {
	w.w = writer
	w.k.reset(key)
	return w
}

func (w *XorWriter) Write(p []byte) (n int, err error) {
	if w.w == nil {
		err = io.ErrClosedPipe
		return
	}
	if w.k.isValid() {
		b := make([]byte, len(p))
		for i, it := range p {
			b[i] = it ^ w.k.next()
		}
		p = b
	}
	n, err = w.w.Write(p)
	return
}

func NewXorWriter(writer io.Writer, key XorKey) *XorWriter {
	return new(XorWriter).Reset(writer, key)
}

package common

import (
	"bytes"
	"io"
	"mksea/crypto"
)

type MetaInfo struct {
	Name  string
	Files []string
}

func (mi *MetaInfo) Append(name string) {
	mi.Files = append(mi.Files, name)
}

func (mi *MetaInfo) Len() int {
	return len(mi.Files)
}

func (mi *MetaInfo) Encode(key crypto.XorKey) (data []byte, err error) {
	defer func() {
		if err != nil {
			err = NewContextError("encoding MetaInfo", err)
		}
	}()
	var buffer bytes.Buffer
	out := crypto.NewXorWriter(&buffer, key)
	errList := NewErrorListCap(1 + len(mi.Files))
	writeBytes := func(b []byte) {
		_, err := out.Write(b)
		errList.Append(err)
	}
	writeInt := func(i int) {
		u := uint64(i)
		var b [8]byte
		for i := range b {
			b[7-i] = byte(u & 0xff)
			u >>= 8
		}
		writeBytes(b[:])
	}
	writeString := func(s string) {
		b := []byte(s)
		writeInt(len(b))
		writeBytes(b)
	}
	writeString(mi.Name)
	writeInt(len(mi.Files))
	for _, it := range mi.Files {
		writeString(it)
	}
	data = buffer.Bytes()
	err = errList.RealError()
	return
}

func (mi *MetaInfo) Decode(data []byte, key crypto.XorKey) (err error) {
	defer func() {
		if err != nil {
			err = NewContextError("decoding MetaInfo", err)
		}
	}()
	in := crypto.NewXorReader(bytes.NewReader(data), key)
	errList := NewErrorListCap(4)
	readBytes := func(b []byte) bool {
		n, err := in.Read(b)
		if err == nil {
			if l := len(b); n < l {
				err = io.ErrShortWrite
			} else if n > l {
				err = io.ErrShortBuffer
			}
		}
		if err != nil {
			errList.Append(err)
			return false
		}
		return true
	}
	readInt := func() (int, bool) {
		var b [8]byte
		if !readBytes(b[:]) {
			return 0, false
		}
		var u uint64
		for _, it := range b {
			u = (u << 8) | uint64(it)
		}
		return int(u), true
	}
	strBuf := make([]byte, 1024)
	readString := func() (string, bool) {
		l, res := readInt()
		if res {
			bl := len(strBuf)
			for l > bl {
				bl *= 2
			}
			if bl != len(strBuf) {
				strBuf = make([]byte, bl)
			}
			res = readBytes(strBuf[:l])
		}
		return string(strBuf[:l]), res
	}
	var hasNext bool
	mi.Name, hasNext = readString()
	if hasNext {
		var count int
		count, hasNext = readInt()
		if hasNext {
			mi.Files = make([]string, count)
			for i := range mi.Files {
				mi.Files[i], hasNext = readString()
				if !hasNext {
					break
				}
			}
		}
	}
	err = errList.RealError()
	return
}

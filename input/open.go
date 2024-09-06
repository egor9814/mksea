package input

import (
	"io"
)

func Open(name string, offset int64) (Interface, error) {
	raw, err := newRawInput(name, offset)
	if err != nil {
		return nil, err
	}
	var rc io.ReadCloser = raw
	if Env.Decode {
		rc = newDecoderInput(rc, rc)
	}
	switch Env.ArchiveFormat {
	case ArchiveNone:
		return newTarInput(rc, rc), nil
	case ArchiveZstd:
		return newZstdInput(rc, rc)
	case ArchiveXz:
		return newXzInput(rc, rc)
	default:
		panic("unreachable")
	}
}

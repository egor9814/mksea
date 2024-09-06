package output

import "github.com/klauspost/compress/zstd"

var Env struct {
	WorkDir            string
	Encode             bool
	EncoderKey         []byte
	ZstdEncoderLevel   zstd.EncoderLevel
	ZstdEncoderThreads int
	XzEncode           bool
}

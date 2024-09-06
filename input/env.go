package input

var Env struct {
	MaxMem        uint64
	ArchiveFormat ArchiveFormat
	Decode        bool
	DecodeKey     []byte
	PasswordTest  []byte
}

func init() {
	Env.MaxMem = 1 * 1024 * 1024 * 1024 // 1GB by default
}

package input

const (
	ArchiveNone ArchiveFormat = iota
	ArchiveZstd
	ArchiveXz
)

func (f ArchiveFormat) Name() string {
	switch f {
	case ArchiveNone:
		return "ArchiveNone"
	case ArchiveZstd:
		return "ArchiveZstd"
	case ArchiveXz:
		return "ArchiveXz"
	default:
		panic("unreachable")
	}
}

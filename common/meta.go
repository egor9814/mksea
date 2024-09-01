package common

type MetaInfo struct {
	Name  string
	Files []string
}

func (i *MetaInfo) Append(name string) {
	i.Files = append(i.Files, name)
}

func (i *MetaInfo) Len() int {
	return len(i.Files)
}

package main

import (
	_ "embed"
	"mksea/common"
)

var metaInfo common.MetaInfo

//go:embed meta.dat
var metaData []byte

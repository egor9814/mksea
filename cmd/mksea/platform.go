package main

import (
	"fmt"
	"runtime"
	"strings"
)

type TargetPlatform uint64

const TargetNone TargetPlatform = 0
const (
	TargetAMD64 TargetPlatform = 1 << iota
	TargetAMD64P32
	Target386
	TargetARM

	TargetWindows
	TargetLinux
	// TODO: test
	// @{
	// TargetDarwin
	// TargetAndroid
	// TargetFreeBSD
	// TargetNetBSD
	// TargetOpenBSD
	// TargetDragonFlyBSD
	// TargetNacl
	// TargetPlan9
	// @}

	TargetArchLast = TargetARM
	TargetOsLast   = TargetLinux

	TargetArchMask = TargetArchLast<<1 - 1
	TargetOsMask   = (TargetOsLast<<1 - 1) ^ TargetArchMask
)

func (p *TargetPlatform) FromString(s string) {
	*p = TargetNone
	slash := strings.Index(s, "/")
	var arch, osname string
	if slash == -1 {
		osname = s
	} else {
		osname = s[:slash]
		arch = s[slash+1:]
	}
	if len(osname) == 0 {
		osname = runtime.GOOS
	}
	if len(arch) == 0 {
		arch = runtime.GOARCH
	}

	for i, it := range []string{"amd64", "amd64p32", "386", "arm"} {
		if it == arch {
			*p |= (TargetPlatform)(int(TargetAMD64) << i)
			break
		}
	}

	for i, it := range []string{"windows", "linux" /* , "darwin", "android", "freebsd", "netbsd", "openbsd", "dragonfly", "nacl", "plan9" */} {
		if it == osname {
			*p |= (TargetPlatform)(int(TargetWindows) << i)
			break
		}
	}
}

func (p TargetPlatform) IsValid() bool {
	return (p&TargetArchMask) != TargetNone && (p&TargetOsMask) != TargetNone
}

func (p TargetPlatform) OsName() string {
	if !p.IsValid() {
		return "<invalid>"
	}
	osnames := map[TargetPlatform]string{
		TargetWindows: "windows",
		TargetLinux:   "linux",
		// TargetDarwin:       "darwin",
		// TargetAndroid:      "android",
		// TargetFreeBSD:      "freebsd",
		// TargetNetBSD:       "netbsd",
		// TargetOpenBSD:      "openbsd",
		// TargetDragonFlyBSD: "dragonfly",
		// TargetNacl:         "nacl",
		// TargetPlan9:        "plan9",
	}
	return osnames[p&TargetOsMask]
}

func (p TargetPlatform) ArchName() string {
	if !p.IsValid() {
		return "<invalid>"
	}
	archnames := map[TargetPlatform]string{
		TargetAMD64:    "amd64",
		TargetAMD64P32: "amd64p32",
		Target386:      "386",
		TargetARM:      "arm",
	}
	return archnames[p&TargetArchMask]
}

func (p TargetPlatform) String() string {
	if !p.IsValid() {
		return "<invalid>"
	}
	return fmt.Sprintf("%s - %s", p.OsName(), p.ArchName())
}

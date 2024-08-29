package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type FileSet map[string]struct{}

func NewFileSet() FileSet {
	return make(FileSet)
}

func (s *FileSet) Resolve(p string, filter func(string, fs.FileInfo) bool) {
	if !filepath.IsAbs(p) {
		p = filepath.Join(workDir, p)
	}
	basepath := filepath.ToSlash(p)
	for _, err := os.Stat(p); err != nil; _, err = os.Stat(p) {
		p = filepath.Dir(p)
	}
	p = filepath.ToSlash(p)
	pattern := wildcardRegexp(basepath)
	filepath.Walk(p, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if filter != nil && !filter(path, info) {
				return nil
			}
			if patternMatched(pattern, filepath.ToSlash(path)) {
				(*s)[path] = struct{}{}
			}
		}
		return nil
	})
}

func wildcardRegexp(pattern string) *regexp.Regexp {
	// based on https://stackoverflow.com/a/74491682
	var result strings.Builder
	result.Grow(len(pattern) * 2)
	result.WriteByte('^')
	for _, it := range pattern {
		switch it {
		case '*':
			result.WriteByte('.')
			result.WriteByte('*')
		case '?':
			result.WriteByte('.')
		case '+', '|', '^', '$', '(', ')', '[', ']', '{', '}':
			result.WriteByte('\\')
			fallthrough
		default:
			result.WriteRune(it)
		}
	}
	result.WriteByte('$')
	return regexp.MustCompile(result.String())
}

func patternMatched(pattern *regexp.Regexp, s string) bool {
	return pattern.Match([]byte(s))
}

func (s *FileSet) Remove(other FileSet) {
	for k := range other {
		delete(*s, k)
	}
}

func (s FileSet) Len() int {
	return len(s)
}

func (s FileSet) List() []string {
	l := make([]string, 0, s.Len())
	for k := range s {
		l = append(l, k)
	}
	return l
}

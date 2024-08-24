package common

import (
	"fmt"
	"strings"
)

type ContextError struct {
	Context string
	Err     error
}

func (err *ContextError) Error() string {
	return fmt.Sprintf("%s: %v", err.Context, err.Err)
}

func NewContextError(context string, err error) *ContextError {
	return &ContextError{
		Context: context,
		Err:     err,
	}
}

type ErrorList []error

func (l *ErrorList) Append(err ...error) {
	for _, it := range err {
		if it != nil {
			if errList, ok := it.(ErrorList); ok {
				*l = append(*l, errList...)
			} else {
				*l = append(*l, it)
			}
		}
	}
}

func (l ErrorList) Len() int {
	return len(l)
}

func (l ErrorList) StringSep(sep string) string {
	s := make([]string, l.Len())
	for i, it := range l {
		s[i] = it.Error()
	}
	return strings.Join(s, sep)
}

func (l ErrorList) String() string {
	return `["` + l.StringSep(`", "`) + `"]`
}

func (l ErrorList) Error() string {
	return l.StringSep("\n")
}

func (l ErrorList) RealError() error {
	if length := l.Len(); length == 0 {
		return nil
	} else if length == 1 {
		return l[0]
	} else {
		return l
	}
}

func NewErrorListCap(count int) ErrorList {
	return make(ErrorList, 0, count)
}

func NewErrorList() ErrorList {
	return NewErrorListCap(2)
}

func NewErrorListFrom(err ...error) ErrorList {
	l := NewErrorListCap(len(err))
	l.Append(err...)
	return l
}

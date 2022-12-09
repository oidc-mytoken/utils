package stringutils

import (
	"fmt"
)

func Wrap(s, wrap string) string {
	return fmt.Sprintf("%s%s%s", wrap, s, wrap)
}

func Unwrap(s, wrap string) string {
	wrapLen := len(wrap)
	sLen := len(s)
	if wrapLen == 0 || sLen == 0 || sLen < 2*wrapLen {
		return s
	}
	if s[0:wrapLen] != wrap || s[sLen-wrapLen:] != wrap {
		return s
	}
	return s[wrapLen : sLen-wrapLen]
}

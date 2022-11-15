package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unsafe"
)

var src rand.Source

func init() {
	src = rand.NewSource(time.Now().UnixNano())
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandASCIIString returns a random string consisting of ASCII characters of the given
// length.
func RandASCIIString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b)) // unsafe is fine here skipcq: GSC-G103
}

// IntersectSlices returns the common elements of two slices
func IntersectSlices(a, b []string) (res []string) {
	for _, bb := range b {
		if StringInSlice(bb, a) {
			res = append(res, bb)
		}
	}
	return
}

// StringInSlice checks if a string is in a slice of strings
func StringInSlice(key string, slice []string) bool {
	for _, s := range slice {
		if s == key {
			return true
		}
	}
	return false
}

// ReplaceStringInSlice replaces all occurrences of a string in a slice with another string
func ReplaceStringInSlice(s *[]string, o, n string, caseSensitive bool) {
	if !caseSensitive {
		o = strings.ToLower(o)
	}
	for i, ss := range *s {
		if !caseSensitive {
			ss = strings.ToLower(ss)
		}
		if o == ss {
			(*s)[i] = n
		}
	}
}

// IsSubSet checks if all strings of a slice 'a' are contained in the slice 'b'
func IsSubSet(a, b []string) bool {
	for _, aa := range a {
		if !StringInSlice(aa, b) {
			return false
		}
	}
	return true
}

// CombineURLPath combines multiple parts of a url
func CombineURLPath(p string, ps ...string) (r string) {
	r = p
	for _, pp := range ps {
		if pp == "" {
			continue
		}
		if r == "" {
			r = pp
			continue
		}
		rAppend := r
		ppAppend := pp
		if strings.HasSuffix(r, "/") {
			rAppend = r[:len(r)-1]
		}
		if strings.HasPrefix(pp, "/") {
			ppAppend = pp[1:]
		}
		r = fmt.Sprintf("%s%c%s", rAppend, '/', ppAppend)
	}
	return
}

// UniqueSlice will remove all duplicates from the given slice of strings
func UniqueSlice(a []string) (unique []string) {
	for _, aa := range a {
		if !StringInSlice(aa, unique) {
			unique = append(unique, aa)
		}
	}
	return
}

// SliceUnion will create a slice of string that contains all strings part of the passed slices
func SliceUnion(a ...[]string) []string {
	res := []string{}
	for _, aa := range a {
		res = append(res, aa...)
	}
	return UniqueSlice(res)
}

// GetTimeIn adds the passed number of seconds to the current time
func GetTimeIn(seconds int64) time.Time {
	return time.Now().Add(time.Duration(seconds) * time.Second)
}

// NewInt64 creates a new *int64
func NewInt64(i int64) *int64 {
	return &i
}

// NewInt creates a new *int
func NewInt(i int) *int {
	return &i
}

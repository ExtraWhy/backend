package util

import (
	"os"
	"unsafe"
)

// B2S converts bytes slice to string without memory allocation.
func B2S(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// S2B converts string to bytes slice without memory allocation.
func S2B(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// ToSlash brings filenames to true slashes
// without superfluous allocations if it possible.
func ToSlash(s string) string {
	var b = S2B(s)
	var bc = b
	var c bool
	for i, v := range b {
		if v == '\\' {
			if !c {
				bc, c = []byte(s), true
			}
			bc[i] = '/'
		}
	}
	return B2S(bc)
}

// ToLower is high performance function to bring filenames to lower case in ASCII
// without superfluous allocations if it possible.
func ToLower(s string) string {
	var b = S2B(s)
	var bc = b
	var c bool
	for i, v := range b {
		if v >= 'A' && v <= 'Z' {
			if !c {
				bc, c = []byte(s), true
			}
			bc[i] |= 0x20
		}
	}
	return B2S(bc)
}

// ToUpper is high performance function to bring filenames to upper case in ASCII
// without superfluous allocations if it possible.
func ToUpper(s string) string {
	var b = S2B(s)
	var bc = b
	var c bool
	for i, v := range b {
		if v >= 'a' && v <= 'z' {
			if !c {
				bc, c = []byte(s), true
			}
			bc[i] &= 0xdf
		}
	}
	return B2S(bc)
}

// ToKey is high performance function to bring filenames to lower case in ASCII
// and true slashes at once without superfluous allocations if it possible.
func ToKey(s string) string {
	var b = S2B(s)
	var bc = b
	var c bool
	for i, v := range b {
		if v >= 'A' && v <= 'Z' {
			if !c {
				bc, c = []byte(s), true
			}
			bc[i] |= 0x20
		} else if v == '\\' {
			if !c {
				bc, c = []byte(s), true
			}
			bc[i] = '/'
		}
	}
	return B2S(bc)
}

// ToID is high performance function to bring filenames to lower case identifier
// with only letters, digits and '_', without superfluous allocations if it possible.
func ToID(s string) string {
	var b = S2B(s)
	var bc = b
	var n int
	for _, v := range b {
		if VarChar[v] || v == '/' {
			n++
		}
	}
	var c bool
	if n != len(b) {
		bc = make([]byte, n)
		c = true
	}
	var i int
	for _, v := range b {
		if VarChar[v] || v == '/' {
			if v >= 'A' && v <= 'Z' {
				if !c {
					bc, c = []byte(s), true
				}
				bc[i] = v | 0x20
			} else if c {
				bc[i] = v
			}
			i++
		}
	}
	return B2S(bc)
}

// JoinPath performs fast join of two UNIX-like path chunks.
func JoinPath(dir, base string) string {
	if dir == "" || dir == "." {
		return base
	}
	if base == "" || base == "." {
		return dir
	}
	if dir[len(dir)-1] == '/' {
		if base[0] == '/' {
			return dir + base[1:]
		} else {
			return dir + base
		}
	}
	if base[0] == '/' {
		return dir + base
	}
	return dir + "/" + base
}

// OS-specific path separator string
const PathSeparator = string(os.PathSeparator)

// JoinFilePath performs fast join of two file path chunks.
// In some cases concatenates with OS-specific separator.
func JoinFilePath(dir, base string) string {
	if dir == "" || dir == "." {
		return base
	}
	if base == "" || base == "." {
		return dir
	}
	var isd = os.IsPathSeparator(dir[len(dir)-1])
	var isb = os.IsPathSeparator(base[0])
	if isd {
		if isb {
			return dir + base[1:]
		} else {
			return dir + base
		}
	}
	if isb {
		return dir + base
	}
	return dir + PathSeparator + base
}

// PathName returns name of file in given file path without extension.
func PathName(fpath string) string {
	var j = len(fpath)
	if j == 0 {
		return ""
	}
	var i = j - 1
	for {
		if fpath[i] == '\\' || fpath[i] == '/' {
			i++
			break
		}
		if fpath[i] == '.' {
			j = i
		}
		if i == 0 {
			break
		}
		i--
	}
	return fpath[i:j]
}

// VarCharFirst is table for fast check that ASCII code is acceptable first symbol of variable.
var VarCharFirst [256]bool = func() (a [256]bool) {
	a['_'] = true
	for c := 'A'; c <= 'Z'; c++ {
		a[c] = true
		a[c+32] = true
	}
	return
}()

// VarChar is table for fast check that ASCII code is acceptable symbol of variable.
var VarChar [256]bool = func() (a [256]bool) {
	a['_'] = true
	for c := 'A'; c <= 'Z'; c++ {
		a[c] = true
		a[c+32] = true
	}
	for c := '0'; c <= '9'; c++ {
		a[c] = true
	}
	return
}()

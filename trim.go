package str

import (
	"unicode"
	"unicode/utf8"
)

// asciiSet is a 32-byte value, where each bit represents the presence of a
// given ASCII character in the set. The 128-bits of the lower 16 bytes,
// starting with the least-significant bit of the lowest word to the
// most-significant bit of the highest word, map to the full range of all
// 128 ASCII characters. The 128-bits of the upper 16 bytes will be zeroed,
// ensuring that any non-ASCII character will be reported as not in the set.
// This allocates a total of 32 bytes even though the upper half
// is unused to avoid bounds checks in asciiSet.contains.
type asciiSet [8]uint32

// contains reports whether c is inside the set.
func (as *asciiSet) contains(c byte) bool {
	return (as[c/32] & (1 << (c % 32))) != 0
}

// makeASCIISet creates a set of ASCII characters and reports whether all
// characters in chars are ASCII.
func makeASCIISet(chars Str) (as asciiSet, ok bool) {
	for _, c := range chars.asBytes() {
		if c >= utf8.RuneSelf {
			return as, false
		}
		as[c/32] |= 1 << (c % 32)
	}
	return as, true
}

// Trim returns a slice of the string s with all leading and
// trailing Unicode code points contained in cutset removed.
func Trim(s, cutset Str) Str {
	if s.Len == 0 || cutset.Len == 0 {
		return s
	}
	if cutset.Len == 1 && cutset.Get(0) < utf8.RuneSelf {
		return trimLeftByte(trimRightByte(s, cutset.Get(0)), cutset.Get(0))
	}
	if as, ok := makeASCIISet(cutset); ok {
		return trimLeftASCII(trimRightASCII(s, &as), &as)
	}
	return trimLeftUnicode(trimRightUnicode(s, cutset), cutset)
}

// TrimLeft returns a slice of the string s with all leading
// Unicode code points contained in cutset removed.
//
// To remove a prefix, use [TrimPrefix] instead.
func TrimLeft(s, cutset Str) Str {
	if s.Len == 0 || cutset.Len == 0 {
		return s
	}
	if cutset.Len == 1 && cutset.Get(0) < utf8.RuneSelf {
		return trimLeftByte(s, cutset.Get(0))
	}
	if as, ok := makeASCIISet(cutset); ok {
		return trimLeftASCII(s, &as)
	}
	return trimLeftUnicode(s, cutset)
}

func trimLeftByte(s Str, c byte) Str {
	for s.Len > 0 && s.Get(0) == c {
		s = s.SliceFrom(1)
	}
	return s
}

func trimLeftASCII(s Str, as *asciiSet) Str {
	for s.Len > 0 {
		if !as.contains(s.Get(0)) {
			break
		}
		s = s.SliceFrom(1)
	}
	return s
}

func trimLeftUnicode(s, cutset Str) Str {
	for s.Len > 0 {
		r, n := rune(s.Get(0)), 1
		if r >= utf8.RuneSelf {
			r, n = utf8.DecodeRuneInString(s.String())
		}
		if !ContainsRune(cutset, r) {
			break
		}
		s = s.SliceFrom(n)
	}
	return s
}

// TrimRight returns a slice of the string s, with all trailing
// Unicode code points contained in cutset removed.
//
// To remove a suffix, use [TrimSuffix] instead.
func TrimRight(s, cutset Str) Str {
	if s.Len == 0 || cutset.Len == 0 {
		return s
	}
	if cutset.Len == 1 && cutset.Get(0) < utf8.RuneSelf {
		return trimRightByte(s, cutset.Get(0))
	}
	if as, ok := makeASCIISet(cutset); ok {
		return trimRightASCII(s, &as)
	}
	return trimRightUnicode(s, cutset)
}

func trimRightByte(s Str, c byte) Str {
	for s.Len > 0 && s.Get(s.Len-1) == c {
		s = s.SliceTo(s.Len - 1)
	}
	return s
}

func trimRightASCII(s Str, as *asciiSet) Str {
	for s.Len > 0 {
		if !as.contains(s.Get(s.Len - 1)) {
			break
		}
		s = s.SliceTo(s.Len - 1)
	}
	return s
}

func trimRightUnicode(s, cutset Str) Str {
	for s.Len > 0 {
		r, n := rune(s.Get(s.Len-1)), 1
		if r >= utf8.RuneSelf {
			r, n = utf8.DecodeLastRuneInString(s.String())
		}
		if !ContainsRune(cutset, r) {
			break
		}
		s = s.SliceTo(s.Len - n)
	}
	return s
}

// indexFunc is the same as IndexFunc except that if
// truth==false, the sense of the predicate function is
// inverted.
func indexFunc(s Str, f func(rune) bool, truth bool) int {
	for i, r := range s.String() {
		if f(r) == truth {
			return i
		}
	}
	return -1
}

// lastIndexFunc is the same as LastIndexFunc except that if
// truth==false, the sense of the predicate function is
// inverted.
func lastIndexFunc(s Str, f func(rune) bool, truth bool) int {
	for i := s.Len; i > 0; {
		r, size := utf8.DecodeLastRuneInString(s.SliceTo(i).String())
		i -= size
		if f(r) == truth {
			return i
		}
	}
	return -1
}

// TrimLeftFunc returns a slice of the string s with all leading
// Unicode code points c satisfying f(c) removed.
func TrimLeftFunc(s Str, f func(rune) bool) Str {
	i := indexFunc(s, f, false)
	if i == -1 {
		return empty
	}
	return s.SliceFrom(i)
}

// TrimRightFunc returns a slice of the string s with all trailing
// Unicode code points c satisfying f(c) removed.
func TrimRightFunc(s Str, f func(rune) bool) Str {
	i := lastIndexFunc(s, f, false)
	if i >= 0 && s.Get(i) >= utf8.RuneSelf {
		_, wid := utf8.DecodeRuneInString(s.SliceFrom(i).String())
		i += wid
	} else {
		i++
	}
	return s.SliceTo(i)
}

// TrimFunc returns a slice of the string s with all leading
// and trailing Unicode code points c satisfying f(c) removed.
func TrimFunc(s Str, f func(rune) bool) Str {
	return TrimRightFunc(TrimLeftFunc(s, f), f)
}

// TrimSpace returns a slice of the string s, with all leading
// and trailing white space removed, as defined by Unicode.
func TrimSpace(s Str) Str {
	// Fast path for ASCII: look for the first ASCII non-space byte
	start := 0
	for ; start < s.Len; start++ {
		c := s.Get(start)
		if c >= utf8.RuneSelf {
			// If we run into a non-ASCII byte, fall back to the
			// slower unicode-aware method on the remaining bytes
			return TrimFunc(s.SliceFrom(start), unicode.IsSpace)
		}
		if asciiSpace[c] == 0 {
			break
		}
	}

	// Now look for the first ASCII non-space byte from the end
	stop := s.Len
	for ; stop > start; stop-- {
		c := s.Get(stop - 1)
		if c >= utf8.RuneSelf {
			// start has been already trimmed above, should trim end only
			return TrimRightFunc(s.Slice(start, stop), unicode.IsSpace)
		}
		if asciiSpace[c] == 0 {
			break
		}
	}

	// At this point s[start:stop] starts and ends with an ASCII
	// non-space bytes, so we're done. Non-ASCII cases have already
	// been handled above.
	return s.Slice(start, stop)
}

// TrimPrefix returns s without the provided leading prefix string.
// If s doesn't start with prefix, s is returned unchanged.
func TrimPrefix(s, prefix Str) Str {
	if !HasPrefix(s, prefix) {
		return s
	}

	return s.SliceFrom(prefix.Len)
}

// TrimSuffix returns s without the provided trailing suffix string.
// If s doesn't end with suffix, s is returned unchanged.
func TrimSuffix(s, suffix Str) Str {
	if !HasSuffix(s, suffix) {
		return s
	}

	return s.SliceTo(s.Len - suffix.Len)
}

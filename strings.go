package str

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// IndexByte returns the index of the first instance of c in s, or -1 if c is not present in s.
func IndexByte(s Str, c byte) int {
	for i, b := range s.asBytes() {
		if b == c {
			return i
		}
	}
	return -1
}

func Equal(s, t Str) bool {
	if s.Len != t.Len {
		return false
	}

	sb := s.asBytes()
	tb := t.asBytes()
	for i, b := range tb {
		if sb[i] != b {
			return false
		}
	}
	return true
}

// Index returns the index of the first instance of substr in s, or -1 if substr is not present in s.
func Index(s, substr Str) int {
	n := substr.Len
	switch {
	case n == 0:
		return 0
	case n == 1:
		return IndexByte(s, substr.Get(0))
	case n == s.Len:
		if substr == s {
			return 0
		}
		return -1
	case n > s.Len:
		return -1
	default:
		// TODO: very inefficient
		for i := 0; i <= s.Len-n; i++ {
			if Equal(s.Slice(i, i+n), substr) {
				return i
			}
		}
		return -1
	}
}

// IndexRune returns the index of the first instance of the Unicode code point
// r, or -1 if rune is not present in s.
// If r is utf8.RuneError, it returns the first instance of any
// invalid UTF-8 byte sequence.
func IndexRune(s Str, r rune) int {
	return strings.IndexRune(s.String(), r)
}

// Contains reports whether substr is within s.
func Contains(s, substr Str) bool {
	return Index(s, substr) >= 0
}

// ContainsRune reports whether the Unicode code point r is within s.
func ContainsRune(s Str, r rune) bool {
	return IndexRune(s, r) >= 0
}

// Cut slices s around the first instance of sep,
// returning the text before and after sep.
// The found result reports whether sep appears in s.
// If sep does not appear in s, cut returns s, "", false.
func Cut(s, sep Str) (before, after Str, found bool) {
	i := Index(s, sep)
	if i < 0 {
		return s, empty, false
	}

	return s.SliceTo(i), s.SliceFrom(i + sep.Len), true
}

// HasPrefix reports whether the string s begins with prefix.
func HasPrefix(s, prefix Str) bool {
	return s.Len >= prefix.Len && Equal(s.SliceTo(prefix.Len), prefix)
}

// HasSuffix reports whether the string s ends with suffix.
func HasSuffix(s, suffix Str) bool {
	return s.Len >= suffix.Len && Equal(s.SliceFrom(s.Len-suffix.Len), suffix)
}

// CutPrefix returns s without the provided leading prefix string
// and reports whether it found the prefix.
// If s doesn't start with prefix, CutPrefix returns s, false.
// If prefix is the empty string, CutPrefix returns s, true.
func CutPrefix(s, prefix Str) (after Str, found bool) {
	if !HasPrefix(s, prefix) {
		return s, false
	}

	return s.SliceFrom(prefix.Len), true
}

// CutSuffix returns s without the provided ending suffix string
// and reports whether it found the suffix.
// If s doesn't end with suffix, CutSuffix returns s, false.
// If suffix is the empty string, CutSuffix returns s, true.
func CutSuffix(s, suffix Str) (before Str, found bool) {
	if !HasSuffix(s, suffix) {
		return s, false
	}

	return s.SliceTo(s.Len - suffix.Len), true
}

// FieldsFunc splits the string s at each run of Unicode code points c satisfying f(c)
// and returns an array of slices of s. If all code points in s satisfy f(c) or the
// string is empty, an empty slice is returned.
//
// FieldsFunc makes no guarantees about the order in which it calls f(c)
// and assumes that f always returns the same value for a given c.
func FieldsFunc(s Str, f func(rune) bool) Iterator {
	return func(yield func(Str) bool) {
		// Find the field start and end indices.
		// Doing this in a separate pass (rather than slicing the string s
		// and collecting the result substrings right away) is significantly
		// more efficient, possibly due to cache effects.
		start := -1 // valid span start if >= 0
		for end, r := range s.String() {
			if f(r) {
				if start >= 0 {
					if !yield(s.Slice(start, end)) {
						return
					}
					// Set start to a negative value.
					// Note: using -1 here consistently and reproducibly
					// slows down this code by a several percent on amd64.
					start = ^start
				}
			} else {
				if start < 0 {
					start = end
				}
			}
		}

		// Last field might end at EOF.
		if start >= 0 {
			if !yield(s.Slice(start, s.Len)) {
				return
			}
		}
	}
}

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

// Fields splits the string s around each instance of one or more consecutive white space
// characters, as defined by unicode.IsSpace, returning a slice of substrings of s or an
// empty slice if s contains only white space.
func Fields(s Str) Iterator {
	// First count the fields.
	// This is an exact count if s is ASCII, otherwise it is an approximation.
	n := 0
	wasSpace := 1
	// setBits is used to track which bits are set in the bytes of s.
	setBits := uint8(0)
	for i := 0; i < s.Len; i++ {
		r := s.Get(i)
		setBits |= r
		isSpace := int(asciiSpace[r])
		n += wasSpace & ^isSpace
		wasSpace = isSpace
	}

	if setBits >= utf8.RuneSelf {
		// Some runes in the input string are not ASCII.
		return FieldsFunc(s, unicode.IsSpace)
	}

	// ASCII fast path
	return func(yield func(Str) bool) {
		fieldStart := 0
		i := 0
		// Skip spaces in the front of the input.
		for i < s.Len && asciiSpace[s.Get(i)] != 0 {
			i++
		}
		fieldStart = i
		for i < s.Len {
			if asciiSpace[s.Get(i)] == 0 {
				i++
				continue
			}
			if !yield(s.Slice(fieldStart, i)) {
				return
			}
			i++
			// Skip spaces in between fields.
			for i < s.Len && asciiSpace[s.Get(i)] != 0 {
				i++
			}
			fieldStart = i
		}
		if fieldStart < s.Len { // Last field might end at EOF.
			yield(s.SliceFrom(fieldStart))
		}
	}
}

// LastIndexByte returns the index of the last instance of c in s, or -1 if c is not present in s.
func LastIndexByte(s Str, c byte) int {
	b := s.asBytes()
	for i := s.Len - 1; i >= 0; i-- {
		if b[i] == c {
			return i
		}
	}
	return -1
}

// Count counts the number of non-overlapping instances of substr in s.
// If substr is an empty string, Count returns 1 + the number of Unicode code points in s.
func Count(s, substr Str) int {
	return strings.Count(s.String(), substr.String())
}

// explode splits s into a slice of UTF-8 strings,
// one string per Unicode character up to a maximum of n (n < 0 means no limit).
// Invalid UTF-8 bytes are sliced individually.
func explode(s Str, n int) Iterator {
	l := utf8.RuneCountInString(s.String())
	if n < 0 || n > l {
		n = l
	}
	return func(yield func(Str) bool) {
		for i := 0; i < n-1; i++ {
			_, size := utf8.DecodeRuneInString(s.String())
			if !yield(s.SliceTo(size)) {
				return
			}
			s = s.SliceFrom(size)
		}
		if n > 0 {
			yield(s)
		}
	}
}

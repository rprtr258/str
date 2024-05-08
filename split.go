package str

// Generic split: splits after each instance of sep,
// including sepSave bytes of sep in the subarrays.
func genSplit(s, sep Str, sepSave, n int) Iterator {
	if n == 0 {
		return EmptyIterator
	}

	if sep.Len == 0 {
		return explode(s, n)
	}

	if n < 0 {
		n = Count(s, sep) + 1
	}

	return func(yield func(Str) bool) {
		n = min(n, s.Len+1)
		n--
		i := 0
		for i < n {
			m := Index(s, sep)
			if m < 0 {
				break
			}
			if !yield(s.SliceTo(m + sepSave)) {
				return
			}
			s = s.SliceFrom(m + sep.Len)
			i++
		}
		yield(s)
	}
}

// SplitN slices s into substrings separated by sep and returns a slice of
// the substrings between those separators.
//
// The count determines the number of substrings to return:
//
//	n > 0: at most n substrings; the last substring will be the unsplit remainder.
//	n == 0: the result is nil (zero substrings)
//	n < 0: all substrings
//
// Edge cases for s and sep (for example, empty strings) are handled
// as described in the documentation for [Split].
//
// To split around the first instance of a separator, see Cut.
func SplitN(s, sep Str, n int) Iterator { return genSplit(s, sep, 0, n) }

// SplitAfterN slices s into substrings after each instance of sep and
// returns a slice of those substrings.
//
// The count determines the number of substrings to return:
//
//	n > 0: at most n substrings; the last substring will be the unsplit remainder.
//	n == 0: the result is nil (zero substrings)
//	n < 0: all substrings
//
// Edge cases for s and sep (for example, empty strings) are handled
// as described in the documentation for SplitAfter.
func SplitAfterN(s, sep Str, n int) Iterator { return genSplit(s, sep, sep.Len, n) }

// Split slices s into all substrings separated by sep and returns a slice of
// the substrings between those separators.
//
// If s does not contain sep and sep is not empty, Split returns a
// slice of length 1 whose only element is s.
//
// If sep is empty, Split splits after each UTF-8 sequence. If both s
// and sep are empty, Split returns an empty slice.
//
// It is equivalent to [SplitN] with a count of -1.
//
// To split around the first instance of a separator, see Cut.
func Split(s, sep Str) Iterator { return genSplit(s, sep, 0, -1) }

// SplitAfter slices s into all substrings after each instance of sep and
// returns a slice of those substrings.
//
// If s does not contain sep and sep is not empty, SplitAfter returns
// a slice of length 1 whose only element is s.
//
// If sep is empty, SplitAfter splits after each UTF-8 sequence. If
// both s and sep are empty, SplitAfter returns an empty slice.
//
// It is equivalent to [SplitAfterN] with a count of -1.
func SplitAfter(s, sep Str) Iterator { return genSplit(s, sep, sep.Len, -1) }

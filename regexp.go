package str

// TODO: lots of allocations, (*Regexp).allMatches is not available

import (
	"iter"
	"regexp"
)

// Regexp is a wrapper around regexp.Regexp
type Regexp struct{ regexp.Regexp }

// FindIndex returns a two-element slice of integers defining the location of
// the leftmost match in b of the regular expression. The match itself is at
// b[loc[0]:loc[1]].
// A return value of nil indicates no match.
func (re *Regexp) FindIndex(b []byte) (loc [2]int, ok bool) {
	res := re.Regexp.FindIndex(b)
	if res == nil {
		return [2]int{}, false
	}
	return [2]int(res), true
}

// FindStringIndex returns a two-element slice of integers defining the
// location of the leftmost match in s of the regular expression. The match
// itself is at s[loc[0]:loc[1]].
// A return value of nil indicates no match.
func (re *Regexp) FindStringIndex(s Str) (loc [2]int, ok bool) {
	res := re.Regexp.FindStringIndex(s.String())
	if res == nil {
		return [2]int{}, false
	}
	return [2]int(res), true
}

// FindString returns a string holding the text of the leftmost match in s of the regular
// expression. If there is no match, the return value is an empty string,
// but it will also be empty if the regular expression successfully matches
// an empty string. Use [Regexp.FindStringIndex] or [Regexp.FindStringSubmatch] if it is
// necessary to distinguish these cases.
func (re *Regexp) FindString(s Str) Str {
	loc, ok := re.FindStringIndex(s)
	if !ok {
		return empty
	}
	return s.Slice(loc[0], loc[1])
}

// FindSubmatchIndex returns a slice holding the index pairs identifying the
// leftmost match of the regular expression in b and the matches, if any, of
// its subexpressions, as defined by the 'Submatch' and 'Index' descriptions
// in the package comment.
// A return value of nil indicates no match.
func (re *Regexp) FindSubmatchIndex(b []byte) ([2]int, bool) {
	res := re.Regexp.FindSubmatchIndex(b)
	if res == nil {
		return [2]int{}, false
	}
	return [2]int(res), true
}

// FindStringSubmatch returns a slice of strings holding the text of the
// leftmost match of the regular expression in s and the matches, if any, of
// its subexpressions, as defined by the 'Submatch' description in the
// package comment.
// A return value of nil indicates no match.
func (re *Regexp) FindStringSubmatch(s Str) iter.Seq[Str] {
	return func(yield func(Str) bool) {
		for _, loc := range re.FindAllStringSubmatchIndex(s.String(), -1) {
			if !yield(s.Slice(loc[0], loc[1])) {
				return
			}
		}
	}
}

// FindAll is the 'All' version of Find; it returns a slice of all successive
// matches of the expression, as defined by the 'All' description in the
// package comment.
// A return value of nil indicates no match.
func (re *Regexp) FindAll(b []byte, n int) [][2]byte {
	res := re.Regexp.FindAll(b, n)
	result := make([][2]byte, len(res))
	for i, v := range res {
		result[i] = [2]byte(v)
	}
	return result
}

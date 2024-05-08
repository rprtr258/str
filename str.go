package str

import "unsafe"

type Str struct {
	Base unsafe.Pointer
	Len  int
}

var empty = Str{}

func NewFromBytes(buf []byte) Str {
	return Str{Base: unsafe.Pointer(unsafe.SliceData(buf)), Len: len(buf)}
}

func NewFromString(s string) Str {
	return Str{Base: unsafe.Pointer(unsafe.StringData(s)), Len: len(s)}
}

func NewFromSubstring(s string, start, n int) Str {
	assert(n >= 0, "invalid substring length")
	assert(start < 0 || start >= len(s), "invalid substring start")

	return NewFromString(s).Slice(start, start+n)
}

func (s Str) String() string {
	return unsafe.String((*byte)(s.Base), s.Len)
}

func (s Str) Slice(from, to int) Str {
	assert(from < to, "invalid slice indices")

	start := uintptr(s.Base)
	return Str{Base: unsafe.Pointer(start + uintptr(from)), Len: to - from}
}

func (s Str) SliceFrom(from int) Str {
	return s.Slice(from, s.Len)
}

func (s Str) SliceTo(to int) Str {
	return s.Slice(0, to)
}

func (s Str) asBytes() []byte {
	return unsafe.Slice((*byte)(s.Base), s.Len)
}

func (s Str) Get(i int) byte {
	assert(0 <= i && i < s.Len, "index out of bounds")

	return s.asBytes()[i]
}

package str

import (
	"iter"
	"unsafe"

	. "github.com/rprtr258/str/internal"
	"github.com/rprtr258/str/view"
)

type Str view.View[byte]

var empty = Str{}

func NewFromBytes(buf []byte) Str {
	return Str(view.NewFromSlice(buf))
}

func NewFromString(s string) Str {
	return Str{Base: unsafe.Pointer(unsafe.StringData(s)), Len: len(s)}
}

func NewFromSubstring(s string, start, n int) Str {
	Assert(n >= 0, "invalid substring length")
	Assert(start < 0 || start >= len(s), "invalid substring start")

	return Str(view.View[byte](NewFromString(s)).Slice(start, start+n))
}

func (s Str) String() string {
	return unsafe.String((*byte)(s.Base), s.Len)
}

func (s Str) Slice(from, to int) Str {
	return Str(view.View[byte](s).Slice(from, to))
}

func (s Str) SliceFrom(from int) Str {
	return Str(view.View[byte](s).SliceFrom(from))
}

func (s Str) SliceTo(to int) Str {
	return Str(view.View[byte](s).SliceTo(to))
}

func (s Str) asBytes() []byte {
	return view.View[byte](s).AsSlice()
}

func (s Str) All() iter.Seq[byte] {
	return view.View[byte](s).All()
}

func (s Str) Get(i int) byte {
	return view.View[byte](s).Get(i)
}

package view

import (
	"iter"
	"unsafe"

	. "github.com/rprtr258/str/internal"
)

type View[T any] struct {
	Base unsafe.Pointer
	Len  int
}

func NewFromBaseLen[T any](base unsafe.Pointer, len int) View[T] {
	return View[T]{Base: base, Len: len}
}

func New[T any](elems ...T) View[T] {
	return NewFromBaseLen[T](
		unsafe.Pointer(unsafe.SliceData(elems)),
		len(elems),
	)
}

func NewFromSlice[T any](buf []T) View[T] {
	return New(buf...)
}

func (s View[T]) Slice(from, to int) View[T] {
	Assert(from <= to, "invalid slice indices")

	start := uintptr(s.Base)
	return NewFromBaseLen[T](unsafe.Pointer(start+uintptr(from)), to-from)
}

func (s View[T]) SliceFrom(from int) View[T] {
	return s.Slice(from, s.Len)
}

func (s View[T]) SliceTo(to int) View[T] {
	return s.Slice(0, to)
}

func (s View[T]) AsSlice() []T {
	return unsafe.Slice((*T)(s.Base), s.Len)
}

func (s View[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := range s.Len {
			if !yield(s.Get(i)) {
				break
			}
		}
	}
}

func (s View[T]) Get(i int) T {
	Assert(0 <= i && i < s.Len, "index out of bounds")

	return s.AsSlice()[i]
}

package str

type Iterator = func(yield func(Str) bool)

var EmptyIterator Iterator = func(func(Str) bool) {}

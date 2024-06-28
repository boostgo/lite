package iterator

type Iterator[T any] struct {
	source       []T
	currentIndex int
}

func New[T any](source []T) *Iterator[T] {
	return &Iterator[T]{
		source: source,
	}
}

func (it *Iterator[T]) Next() bool {
	return it.currentIndex < len(it.source)
}

func (it *Iterator[T]) Get() (T, bool) {
	if it.currentIndex >= len(it.source) {
		return *new(T), false
	}

	defer func() {
		it.currentIndex++
	}()
	return it.source[it.currentIndex], true
}

func (it *Iterator[T]) TryGet() T {
	value, ok := it.Get()
	if !ok {
		return *new(T)
	}

	return value
}

func (it *Iterator[T]) Skip(count int) *Iterator[T] {
	for i := 0; i < count; i++ {
		_ = it.TryGet()
	}
	return it
}

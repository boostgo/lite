package chunk

type Chunk[T ~[]E, E any] struct {
	chunks   []T
	page     int
	maxPages int
}

func New[T ~[]E, E any](source T, chunkSize int) *Chunk[T, E] {
	return &Chunk[T, E]{
		chunks:   divide(source, chunkSize, len(source)/chunkSize),
		page:     0,
		maxPages: len(source) / chunkSize,
	}
}

func (chunk *Chunk[T, E]) Get() []T {
	return chunk.chunks
}

func (chunk *Chunk[T, E]) Next() bool {
	return chunk.page < chunk.maxPages
}

func (chunk *Chunk[T, E]) Read() T {
	if chunk.page >= chunk.maxPages {
		return nil
	}

	defer func() {
		chunk.page++
	}()
	return chunk.chunks[chunk.page]
}

func divide[T ~[]E, E any](source T, size, length int) []T {
	chunks := make([]T, 0, length)

	for i := 0; i < len(source); i += size {
		end := i + size
		if end > len(source) {
			end = len(source)
		}

		chunks = append(chunks, source[i:end])
	}

	return chunks
}

package chunk

type Chunk[T ~[]E, E any] struct {
	source    T
	chunkSize int

	chunksLength int
}

func New[T ~[]E, E any](source T, chunkSize int) *Chunk[T, E] {
	return &Chunk[T, E]{
		source:    source,
		chunkSize: chunkSize,

		chunksLength: len(source) / chunkSize,
	}
}

func (chunk *Chunk[T, E]) Divide() []T {
	chunks := make([]T, 0, chunk.chunksLength)

	for i := 0; i < len(chunk.source); i += chunk.chunkSize {
		end := i + chunk.chunkSize
		if end > len(chunk.source) {
			end = len(chunk.source)
		}

		chunks = append(chunks, chunk.source[i:end])
	}

	return chunks
}

package list

func Chunk[T ~[]E, E any](source T, size int) []T {
	if size <= 0 {
		size = len(source)
	}

	chunks := make([]T, 0, len(source)/size+1)

	for i := 0; i < len(source); i += size {
		end := i + size
		if end > len(source) {
			end = len(source)
		}

		chunks = append(chunks, source[i:end])
	}

	return chunks
}

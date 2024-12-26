package list

// Chunk divide slice for sub-slices by provided chunk size.
// Example:
//
//	texts := []string{"text #1", "text #2", "text #3", "text #4", "text #5"}
//	chunks := list.Chunk(texts, 2)
//	fmt.Println(chunks) // [[text #1 text #2] [text #3 text #4] [text #5]]
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

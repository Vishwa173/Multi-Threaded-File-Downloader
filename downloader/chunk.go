package downloader

import "fmt"

func SplitIntoChunks(size int64, chunkCount int) ([]Chunk, error) {
	if size <= 0 {
		return nil, fmt.Errorf("size must be greater than zero")
	}
	if chunkCount <= 0 {
		return nil, fmt.Errorf("chunkCount must be greater than zero")
	}

	chunkSize := size / int64(chunkCount)

	chunks := make([]Chunk, chunkCount)
	var start int64 = 0

	for i := 0; i < chunkCount; i++ {
		end := start + chunkSize - 1
		if i == chunkCount-1 {
			end = size - 1 
		}
		chunks[i] = Chunk{
			Index: i,
			Start: start,
			End: end,
		}
		start = end + 1
	}

	return chunks, nil
}
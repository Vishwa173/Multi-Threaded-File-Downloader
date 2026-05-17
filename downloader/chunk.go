package downloader

import "fmt"

const DefaultChunkSize = 8 * 1024 * 1024

func SplitIntoChunks(
	size int64,
	chunkSize int64,
) ([]Chunk, error) {

	if size <= 0 {
		return nil, fmt.Errorf(
			"invalid file size",
		)
	}

	if chunkSize <= 0 {
		return nil, fmt.Errorf(
			"invalid chunk size",
		)
	}

	var chunks []Chunk

	var start int64 = 0
	index := 0

	for start < size {

		end := start + chunkSize - 1

		if end >= size {
			end = size - 1
		}

		chunks = append(chunks, Chunk{
			Index:   index,
			Start:   start,
			End:     end,
			Status:  ChunkPending,
			Retries: 0,
		})

		start = end + 1
		index++
	}

	return chunks, nil
}
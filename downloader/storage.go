package downloader

import (
	"fmt"
	"os"
	"path/filepath"
)

func ChunkExists(
	tempDir string,
	chunk Chunk,
) bool {

	partPath := filepath.Join(
		tempDir,
		fmt.Sprintf(
			"chunk-%d.part",
			chunk.Index,
		),
	)

	info, err := os.Stat(partPath)

	if err != nil {
		return false
	}

	expectedSize :=
		chunk.End - chunk.Start + 1

	return info.Size() == expectedSize
}
package downloader

import(
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
	"multi-threaded-downloader/utils"
)

const maxRetries = 3

func DownloadChunk(url string, chunk Chunk, tempDir string) error {

    partPath := filepath.Join(tempDir, fmt.Sprintf("chunk-%d.part", chunk.Index))

    f, err := os.Create(partPath)
    if err != nil {
        return fmt.Errorf("chunk %d: could not create part file: %w", chunk.Index, err)
    }
    defer f.Close()

    for attempt := 1; attempt <= maxRetries; attempt++ {

        body, err := utils.DownloadRange(url, chunk.Start, chunk.End)
        if err != nil {
            if attempt == maxRetries {
                return fmt.Errorf("chunk %d: failed after %d attempts: %w", chunk.Index, attempt, err)
            }

            sleep := time.Duration(attempt*attempt) * 200 * time.Millisecond
            time.Sleep(sleep)
            continue
        }

        _, copyErr := io.Copy(f, body)
        body.Close()

        if copyErr != nil {
            if attempt == maxRetries {
                return fmt.Errorf("chunk %d: copy failed after %d attempts: %w", chunk.Index, attempt, copyErr)
            }

            sleep := time.Duration(attempt*attempt) * 200 * time.Millisecond
            time.Sleep(sleep)
            continue
        }

        return nil
    }

    return fmt.Errorf("chunk %d: unreachable error", chunk.Index)
}

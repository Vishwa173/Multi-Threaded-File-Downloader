package downloader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"multi-threaded-downloader/utils"
)

const maxRetries = 3

func Worker(
	workerID int,
	url string,
	tempDir string,
	scheduler *Scheduler,
	metrics *Metrics,
	errCh chan error,
) {
	LogWorker(workerID, "started")
	chunkIndex := 0

	for {
		metrics.mu.RLock()
		workerSpeed := metrics.WorkerMetrics[workerID].AverageSpeed

		metrics.mu.RUnlock()
		chunk, ok := scheduler.GetNextChunk(
			workerSpeed,
			chunkIndex,
		)

		if !ok {
			break
		}

		chunk.WorkerID = workerID
		chunk.Status = ChunkDownloading

		startTime := time.Now()

		err := DownloadChunk(
			url,
			chunk,
			tempDir,
			metrics,
		)

		if err != nil {
			metrics.MarkChunkFailed()
			errCh <- err
			continue
		}

		duration := time.Since(startTime).Seconds()
		bytes := chunk.End - chunk.Start + 1
		speed := float64(bytes) / duration

		metrics.UpdateWorkerSpeed(
			workerID,
			speed,
		)

		metrics.MarkChunkCompleted()

		LogChunk(
			chunk.Index,
			fmt.Sprintf(
				"worker=%d speed=%.2f MB/s",
				workerID,
				speed/(1024*1024),
			),
		)

		chunkIndex++
	}
	LogWorker(workerID, "finished")
}

func DownloadChunk(
	url string,
	chunk Chunk,
	tempDir string,
	metrics *Metrics,
) error {

	partPath := filepath.Join(
		tempDir,
		fmt.Sprintf("chunk-%d.part", chunk.Index),
	)

	f, err := os.Create(partPath)
	if err != nil {
		return fmt.Errorf(
			"chunk %d: create failed: %w",
			chunk.Index,
			err,
		)
	}

	defer f.Close()

	for attempt := 1; attempt <= maxRetries; attempt++ {

		body, err := utils.DownloadRange(
			url,
			chunk.Start,
			chunk.End,
		)

		if err != nil {

			if attempt == maxRetries {
				return fmt.Errorf(
					"chunk %d failed after retries: %w",
					chunk.Index,
					err,
				)
			}

			sleep := time.Duration(
				attempt*attempt,
			) * 300 * time.Millisecond

			time.Sleep(sleep)

			continue
		}

		written, copyErr := io.Copy(f, body)

		body.Close()

		if copyErr != nil {

			if attempt == maxRetries {
				return fmt.Errorf(
					"chunk %d copy failed: %w",
					chunk.Index,
					copyErr,
				)
			}

			sleep := time.Duration(
				attempt*attempt,
			) * 300 * time.Millisecond

			time.Sleep(sleep)

			continue
		}

		metrics.AddDownloadedBytes(written)
		metrics.UpdateWorker(chunk.WorkerID, written)

		return nil
	}

	return fmt.Errorf(
		"chunk %d unreachable error",
		chunk.Index,
	)
}
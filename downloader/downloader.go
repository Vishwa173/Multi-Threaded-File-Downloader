package downloader

import(
	"fmt"
	"os"
	"net/http"
	"sync"
	"strconv"
	"time"
	"multi-threaded-downloader/utils"
	"strings"
)

func getFileSize(rawURL string) (int64, error) {
    finalURL, err := resolveFinalURL(rawURL)
    if err != nil {
        return 0, fmt.Errorf("failed to resolve URL: %w", err)
    }

    headResp, err := http.Head(finalURL)
    if err == nil && headResp.StatusCode == http.StatusOK {
        sizeStr := headResp.Header.Get("Content-Length")
        if sizeStr != "" {
            size, err := strconv.ParseInt(sizeStr, 10, 64)
            if err == nil && size > 0 {
                return size, nil
            }
        }
    }

    req, err := http.NewRequest("GET", finalURL, nil)
    if err != nil {
        return 0, fmt.Errorf("failed to create request: %w", err)
    }

    req.Header.Set("Range", "bytes=0-1")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return 0, fmt.Errorf("range request failed: %w", err)
    }
    defer resp.Body.Close()

    cr := resp.Header.Get("Content-Range")
    if cr == "" {
        return 0, fmt.Errorf("server does not support range requests")
    }

    parts := strings.Split(cr, "/")
    if len(parts) != 2 {
        return 0, fmt.Errorf("invalid Content-Range format: %s", cr)
    }

    sizeStr := parts[1]
    size, err := strconv.ParseInt(sizeStr, 10, 64)
    if err != nil {
        return 0, fmt.Errorf("failed to parse file size: %w", err)
    }

    return size, nil
}


func resolveFinalURL(url string) (string, error) {
    client := &http.Client{
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            if len(via) >= 10 {
                return fmt.Errorf("too many redirects")
            }
            return nil
        },
    }

    resp, err := client.Head(url)
    if err == nil {
        return resp.Request.URL.String(), nil
    }

    resp, err = client.Get(url)
    if err != nil {
        return "", err
    }
    resp.Body.Close()

    return resp.Request.URL.String(), nil
}

func DownloadFile(url, output string, workers int, expectedHash string,) error {
	InitLogger()

	size, err := getFileSize(url)
	if err != nil {
		return fmt.Errorf(
			"error getting file size: %w",
			err,
		)
	}

	checkpointPath := output + ".meta.json"

	checkpoint, err := LoadCheckpoint(
		checkpointPath,
	)

	if err == nil {
		LogWorker(
			-1,
			"restoring checkpoint",
		)
	}

	tempDir := output + ".parts"
	err = os.MkdirAll(
		tempDir,
		0755,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to create parts dir: %w",
			err,
		)
	}

	sources := []*Source{
		{
			URL:     url,
			Healthy: true,
		},
	}

	scheduler := NewScheduler(
		size,
		DefaultChunkSize,
		sources,
	)

	if checkpoint != nil {
		scheduler.NextByte = checkpoint.NextByte
		scheduler.FailedQueue = checkpoint.FailedChunks
		scheduler.CompletedChunks = checkpoint.CompletedChunks
	}

	metrics := NewMetrics(size)
	var wg sync.WaitGroup
	errCh := make(chan error, workers)
	done := make(chan struct{})

	go func() {
		ticker := time.NewTicker(
			3 * time.Second,
		)

		defer ticker.Stop()

		for {
			select {
			case <-done:
				SaveCheckpoint(
					checkpointPath,
					scheduler,
					scheduler.CompletedChunks,
				)
				return

			case <-ticker.C:
				SaveCheckpoint(
					checkpointPath,
					scheduler,
					scheduler.CompletedChunks,
				)
			}
		}
	}()

	for i := 0; i < workers; i++ {
		metrics.RegisterWorker(i)
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			Worker(
				workerID,
				url,
				tempDir,
				scheduler,
				metrics,
				errCh,
			)

		}(i)
	}

	go func() {
		ticker := time.NewTicker(
			500 * time.Millisecond,
		)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				metrics.mu.RLock()

				downloaded := metrics.DownloadedBytes
				total := metrics.TotalBytes
				metrics.mu.RUnlock()

				percent :=
					float64(downloaded) /
						float64(total) * 100

				printProgress(
					percent,
					downloaded,
					total,
				)

			case <-done:
				return
			}
		}
	}()

	wg.Wait()
	close(done)
	close(errCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	err = utils.MergeChunks(
		output,
		workers,
		tempDir,
	)

	if err != nil {
		return fmt.Errorf(
			"merge failed: %w",
			err,
		)
	}

	if expectedHash != "" {
		LogWorker(-1,"verifying sha256 integrity",)
		err = ValidateSHA256(output,expectedHash,)

		if err != nil {
			return fmt.Errorf(
				"integrity validation failed: %w",
				err,
			)
		}

		LogWorker(-1,"integrity verification passed",)
	}

	hash, err := ComputeSHA256(output)
	if err == nil {
		LogWorker(-1,fmt.Sprintf("final sha256: %s",hash,),
		)
	}

	os.RemoveAll(tempDir)
	os.Remove(checkpointPath)

	LogWorker(-1,"download completed",)

	return nil
}
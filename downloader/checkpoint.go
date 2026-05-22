package downloader

import (
	"encoding/json"
	"os"
)

type Checkpoint struct {
	NextByte int64
	CompletedChunks []Chunk
	FailedChunks []Chunk
}

func SaveCheckpoint(
	path string,
	scheduler *Scheduler,
	completed []Chunk,
) error {

	scheduler.mu.Lock()

	checkpoint := Checkpoint{
		NextByte:        scheduler.NextByte,
		CompletedChunks: completed,
		FailedChunks:    scheduler.FailedQueue,
	}

	scheduler.mu.Unlock()

	data, err := json.MarshalIndent(
		checkpoint,
		"",
		"  ",
	)

	if err != nil {
		return err
	}

	return os.WriteFile(
		path,
		data,
		0644,
	)
}

func LoadCheckpoint(
	path string,
) (*Checkpoint, error) {

	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var checkpoint Checkpoint

	err = json.Unmarshal(
		data,
		&checkpoint,
	)

	if err != nil {
		return nil, err
	}

	return &checkpoint, nil
}
package downloader

import "time"

type ChunkStatus string

const (
	ChunkPending     ChunkStatus = "PENDING"
	ChunkDownloading ChunkStatus = "DOWNLOADING"
	ChunkCompleted   ChunkStatus = "COMPLETED"
	ChunkFailed      ChunkStatus = "FAILED"
)

type Chunk struct {
	Index    int
	Start    int64
	End      int64
	Status   ChunkStatus
	Retries  int
	WorkerID int
}

type WorkerStats struct {
	WorkerID        int
	ChunksHandled   int
	BytesDownloaded int64
	LastActivity    time.Time
	Failures        int
	CurrentSpeed float64
	AverageSpeed float64
}

type Source struct {
	URL string
	Failures int
	AverageSpeed float64
	ActiveWorkers int
	Healthy bool
}
package downloader

import (
	"sync"
	"time"
)

type Metrics struct {
	mu sync.RWMutex

	StartTime time.Time

	TotalBytes      int64
	DownloadedBytes int64

	CompletedChunks int
	FailedChunks    int

	WorkerMetrics map[int]*WorkerStats
}

func NewMetrics(totalBytes int64) *Metrics {
	return &Metrics{
		StartTime:     time.Now(),
		TotalBytes:    totalBytes,
		WorkerMetrics: make(map[int]*WorkerStats),
	}
}

func (m *Metrics) AddDownloadedBytes(n int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.DownloadedBytes += n
}

func (m *Metrics) RegisterWorker(workerID int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.WorkerMetrics[workerID] = &WorkerStats{
		WorkerID:     workerID,
		LastActivity: time.Now(),
	}
}

func (m *Metrics) UpdateWorker(workerID int, bytes int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	worker := m.WorkerMetrics[workerID]

	worker.BytesDownloaded += bytes
	worker.ChunksHandled++
	worker.LastActivity = time.Now()
}

func (m *Metrics) MarkChunkCompleted() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.CompletedChunks++
}

func (m *Metrics) MarkChunkFailed() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.FailedChunks++
}
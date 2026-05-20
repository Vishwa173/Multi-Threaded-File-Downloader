package downloader

import (
	"sync"
	"time"
)

type Scheduler struct {
	mu sync.Mutex

	NextByte int64
	FileSize int64

	BaseChunkSize int64

	FailedQueue []Chunk

	Sources []*Source
}

func NewScheduler(
	fileSize int64,
	baseChunkSize int64,
	sources []*Source,
) *Scheduler {

	return &Scheduler{
		NextByte:      0,
		FileSize:      fileSize,
		BaseChunkSize: baseChunkSize,
		Sources:       sources,
	}
}

func (s *Scheduler) GetNextChunk(
	workerSpeed float64,
	index int,
) (Chunk, bool) {

	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.FailedQueue) > 0 {

		chunk := s.FailedQueue[0]

		s.FailedQueue =
			s.FailedQueue[1:]

		return chunk, true
	}

	if s.NextByte >= s.FileSize {
		return Chunk{}, false
	}

	chunkSize := s.BaseChunkSize
	multiplier := workerSpeed /
		(10 * 1024 * 1024)

	if multiplier < 0.5 {
		multiplier = 0.5
	}

	if multiplier > 4 {
		multiplier = 4
	}

	chunkSize = max(
		int64(float64(chunkSize)*multiplier),
		1*1024*1024,
	)

	start := s.NextByte

	end := start + chunkSize - 1

	if end >= s.FileSize {
		end = s.FileSize - 1
	}

	s.NextByte = end + 1

	return Chunk{
		Index:   index,
		Start:   start,
		End:     end,
		Status:  ChunkPending,
		Retries: 0,
	}, true
}

func (s *Scheduler) RequeueChunk(
	chunk Chunk,
) {

	chunk.Status = ChunkFailed
	chunk.Retries++

	if chunk.Retries >= maxRetries {
		return
	}

	backoff := time.Duration(
		chunk.Retries*chunk.Retries,
	) * 500 * time.Millisecond

	time.Sleep(backoff)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.FailedQueue = append(
		s.FailedQueue,
		chunk,
	)
}

func (s *Scheduler) SelectBestSource() *Source {
	s.mu.Lock()
	defer s.mu.Unlock()

	var best *Source
	for _, source := range s.Sources {
		if !source.Healthy {
			continue
		}
		if best == nil {
			best = source
			continue
		}

		bestScore := best.AverageSpeed / float64(best.ActiveWorkers+1)

		currentScore := source.AverageSpeed / float64(source.ActiveWorkers+1)

		if currentScore > bestScore {
			best = source
		}
	}

	if best == nil {
		return s.Sources[0]
	}

	best.ActiveWorkers++
	return best
}

func (s *Scheduler) ReleaseSource(
	source *Source,
	speed float64,
	success bool,
) {
	s.mu.Lock()
	defer s.mu.Unlock()
	source.ActiveWorkers--

	if success {
		alpha := 0.3

		if source.AverageSpeed == 0 {
			source.AverageSpeed = speed
		} else {
			source.AverageSpeed = alpha*speed + (1-alpha)*source.AverageSpeed
		}

		source.Failures = 0
		source.Healthy = true
	} else {
		source.Failures++

		if source.Failures >= 3 {
			source.Healthy = false
		}
	}
}
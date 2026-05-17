package downloader

import "sync"

type Scheduler struct {
	mu sync.Mutex

	NextByte int64
	FileSize int64

	BaseChunkSize int64
}

func NewScheduler(
	fileSize int64,
	baseChunkSize int64,
) *Scheduler {

	return &Scheduler{
		NextByte:      0,
		FileSize:      fileSize,
		BaseChunkSize: baseChunkSize,
	}
}

func (s *Scheduler) GetNextChunk(
	workerSpeed float64,
	index int,
) (Chunk, bool) {

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.NextByte >= s.FileSize {
		return Chunk{}, false
	}

	chunkSize := s.BaseChunkSize

	multiplier := workerSpeed / (10 * 1024 * 1024)
	
	if multiplier < 0.5 {
		multiplier = 0.5
	}
	
	if multiplier > 4 {
		multiplier = 4
	}

	chunkSize = max(int64(float64(chunkSize) * multiplier,), 1 * 1024 * 1024)

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
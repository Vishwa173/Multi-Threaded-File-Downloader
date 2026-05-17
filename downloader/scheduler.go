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

	if workerSpeed > 30*1024*1024 {
		chunkSize *= 2
	}

	if workerSpeed < 10*1024*1024 {
		chunkSize /= 2
	}

	if chunkSize < 1*1024*1024 {
		chunkSize = 1 * 1024 * 1024
	}

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
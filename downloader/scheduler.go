package downloader

type Scheduler struct {
	ChunkQueue chan Chunk
}

func NewScheduler(chunks []Chunk) *Scheduler {
	queue := make(chan Chunk, len(chunks))

	for _, chunk := range chunks {
		queue <- chunk
	}

	close(queue)

	return &Scheduler{
		ChunkQueue: queue,
	}
}
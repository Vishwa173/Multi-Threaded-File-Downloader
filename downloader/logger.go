package downloader

import (
	"log"
	"os"
)

var Logger *log.Logger

func InitLogger() {
	Logger = log.New(
		os.Stdout,
		"[DOWNLOADER] ",
		log.LstdFlags,
	)
}

func LogWorker(workerID int, msg string) {
	Logger.Printf("[WORKER %d] %s", workerID, msg)
}

func LogChunk(chunkID int, msg string) {
	Logger.Printf("[CHUNK %d] %s", chunkID, msg)
}

func LogError(err error) {
	Logger.Printf("[ERROR] %v", err)
}
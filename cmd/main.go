package main

import (
	"fmt"
	"os"
	"strconv"

	"multi-threaded-downloader/downloader"
)

func main() {

	if len(os.Args) < 4 {
		fmt.Println(
			"Usage: go run ./cmd/main.go <url> <output-file> <threads> [sha256]",
		)
		return
	}

	url := os.Args[1]
	output := os.Args[2]
	threadsStr := os.Args[3]
	var expectedHash string

	if len(os.Args) >= 5 {
		expectedHash = os.Args[4]
	}

	threads, err := strconv.Atoi(
		threadsStr,
	)

	if err != nil {

		fmt.Println(
			"Invalid thread count:",
			err,
		)

		return
	}

	err = downloader.DownloadFile(
		url,
		output,
		threads,
		expectedHash,
	)

	if err != nil {

		fmt.Println(
			"Download failed:",
			err,
		)

		return
	}

	fmt.Println(
		"Download completed:",
		output,
	)
}
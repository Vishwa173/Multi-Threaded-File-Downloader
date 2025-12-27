package main

import (
    "fmt"
    "os"
    "strconv"

    "multi-threaded-downloader/downloader"
)

func main() {
    if len(os.Args) < 4 {
        fmt.Println("Usage: go run ./cmd/main.go <url> <output-file> <threads>")
        return
    }

    url := os.Args[1]
    output := os.Args[2]
    threadsStr := os.Args[3]

    threads, err := strconv.Atoi(threadsStr)
    if err != nil {
        fmt.Println("Invalid thread count:", err)
        return
    }

    err = downloader.DownloadFile(url, output, threads)
    if err != nil {
        fmt.Println("Download failed:", err)
        return
    }

    fmt.Println("Download completed:", output)
}


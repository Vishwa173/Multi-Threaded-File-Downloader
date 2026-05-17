package downloader

import(
	"fmt"
	"strings"
)

func progressBar(percent float64)string{
	filled := int(percent / 2) 
    return strings.Repeat("=", filled) + strings.Repeat("-", 50-filled)
}

func printProgress(
	percent float64,
	downloaded int64,
	total int64,
) {

	fmt.Printf(
		"\rProgress: %.2f%% (%d/%d MB)",
		percent,
		downloaded/(1024*1024),
		total/(1024*1024),
	)
}
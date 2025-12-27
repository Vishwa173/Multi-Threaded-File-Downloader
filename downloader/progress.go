package downloader

import(
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func getDownloadedBytes(tempDir string, chunkCount int)int64{
	var total int64 = 0

	for i := 0; i < chunkCount; i++ {
		partPath := filepath.Join(tempDir,fmt.Sprintf("chunk-%d.part",i))
		info,err := os.Stat(partPath)
		if err == nil{
			total += info.Size()
		}
	}

	return total
}

func progressBar(percent float64)string{
	filled := int(percent / 2) 
    return strings.Repeat("=", filled) + strings.Repeat("-", 50-filled)
}

func printProgress(percent float64, speed float64, downloaded, total int64){
	fmt.Printf("\r[%-50s] %.2f%%  (%.2f MB / %.2f MB)  %.2f MB/s",
        progressBar(percent),
        percent,
        float64(downloaded)/1e6,
        float64(total)/1e6,
        speed/1e6,
    )
	os.Stdout.Sync()
}


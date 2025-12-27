package utils

import(
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func CreateTempDir(base string)(string,error){
	if base == ""{
		base = os.TempDir()
	}

	dirName := fmt.Sprintf("downloader-%d",time.Now().UnixNano())
	fullPath := filepath.Join(base,dirName)

	err := os.MkdirAll(fullPath,0755)
	if err != nil{
		return "", fmt.Errorf("error creating dir: %w",err)
	}

	return fullPath,nil
}

func MergeChunks(outputFile string, chunkCount int, tempDir string)error{
	f,err := os.Create(outputFile)
	if err != nil{
		return fmt.Errorf("error creating output file: %w",err)
	}
	defer f.Close()

	for i := 0; i < chunkCount; i++{
		partPath := filepath.Join(tempDir,fmt.Sprintf("chunk-%d.part",i))

		partFile,err := os.Open(partPath)
		if err != nil{
			return fmt.Errorf("error in opening part file %d: %w",i,err)
		}

		_, err = io.Copy(f,partFile)
		partFile.Close()
		if err != nil{
			return fmt.Errorf("error in copying part file %d: %w",i,err)
		}

		err = os.Remove(partPath)
		if err != nil{
			return fmt.Errorf("failed to delete partfile %d: %w",i,err)
		}
	}
	
	return nil;
}
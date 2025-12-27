package utils

import (
	"fmt"
	"io"
	"net/http"
)

func DownloadRange(url string, start, end int64) (io.ReadCloser,error){
	req,err := http.NewRequest("GET",url,nil)
	if err != nil{
		return nil, fmt.Errorf("Error creating request : %v",err)
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d",start,end))

	client := http.Client{}
	resp,err := client.Do(req)
	if err != nil{
		return nil, fmt.Errorf("Error creating request : %v", err)
	}

	if resp.StatusCode != http.StatusPartialContent{
		return nil, fmt.Errorf("Expected 206 partial content, got %d", resp.StatusCode)
	}

	return resp.Body, nil
}
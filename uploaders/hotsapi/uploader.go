package hotsapi

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

const (
	endpoint = "http://hotsapi.net/api/v1/replays"
)

var (
	client = &http.Client{
		Timeout: 15 * time.Second,
	}
)

type Uploader struct{}

type UploadResult struct {
	Success      bool   `json:"success"`
	Status       string `json:"status"`
	ID           int    `json:"id"`
	OriginalName string `json:"originalName"`
	Filename     string `json:"filename"`
	Url          string `json:"url"`
}

func NewUploader() *Uploader {
	return &Uploader{}
}

func uploadReplay(path string) (result UploadResult, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "replay.StormReplay")
	if err != nil {
		return
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return
	}

	return
}

func (u *Uploader) UploadReplay(replay string) (result string, err error) {
	res, err := uploadReplay(replay)
	if err != nil {
		return
	}

	return res.Status, nil
}

func (u *Uploader) UploadReplays(replays []string) ([]string, []error) {
	results := make([]string, len(replays))
	errs := make([]error, len(replays))

	for i, path := range replays {
		result, err := uploadReplay(path)
		if err != nil {
			errs[i] = err
		} else {
			results[i] = result.Status
		}
	}

	return results, errs
}

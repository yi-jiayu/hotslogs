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
	endpoint = "https://hotsapi.net/api/v1/replays"
)

type Uploader struct {
	Client *http.Client
}

type Response struct {
	Success      bool   `json:"success"`
	Status       string `json:"status"`
	ID           int    `json:"id"`
	OriginalName string `json:"originalName"`
	Filename     string `json:"filename"`
	Url          string `json:"url"`
}

func NewUploader() *Uploader {
	return &Uploader{
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (u *Uploader) uploadReplay(path string) (result Response, err error) {
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

	resp, err := u.Client.Do(req)
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

func (u *Uploader) UploadReplay(replay string) (result Response, err error) {
	res, err := u.uploadReplay(replay)
	if err != nil {
		return
	}

	return res, nil
}

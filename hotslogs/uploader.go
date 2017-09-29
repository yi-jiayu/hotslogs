package hotslogs

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/satori/go.uuid"
)

const (
	accessKeyID     = "AKIAIESBHEUH4KAAG4UA"
	secretAccessKey = "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"
	s3BucketName    = "heroesreplays"
	s3BucketRegion  = "us-west-2"

	endpoint = "https://www.hotslogs.com/UploadFile?Source="
)

var (
	staticCredentials = credentials.NewStaticCredentials(accessKeyID, secretAccessKey, "")
)

type Uploader struct {
	uploader *s3manager.Uploader
}

func NewUploader() *Uploader {
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: staticCredentials,
		Region:      aws.String(s3BucketRegion),
	}))
	uploader := s3manager.NewUploader(sess)

	return &Uploader{
		uploader: uploader,
	}
}

func (u *Uploader) uploadReplayToS3(key string, file io.Reader) error {
	params := &s3manager.UploadInput{
		Bucket: aws.String(s3BucketName),
		Key:    &key,
		Body:   file,
	}

	_, err := u.uploader.Upload(params)
	return err
}

func (u *Uploader) getUploadResult(key string) (string, error) {
	URL := endpoint + "&FileName=" + key
	resp, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func (u *Uploader) UploadReplay(replay string) (string, error) {
	UUID := uuid.NewV4()
	key := fmt.Sprintf("%s.StormReplay", UUID)

	file, err := os.Open(replay)
	if err != nil {
		return "", err
	}

	err = u.uploadReplayToS3(key, file)
	if err != nil {
		return "", err
	}

	result, err := u.getUploadResult(key)
	return result, err
}

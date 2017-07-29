package hotslogs

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/satori/go.uuid"
)

const (
	AccessKeyID     = "AKIAIESBHEUH4KAAG4UA"
	SecretAccessKey = "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"
	S3BucketName    = "heroesreplays"
	S3BucketRegion  = "us-west-2"

	UploadEndpoint = "https://www.hotslogs.com/UploadFile?Source="
)

var (
	StaticCredentials = credentials.NewStaticCredentials(AccessKeyID, SecretAccessKey, "")
)

func ListNewReplays(replayDir string, since time.Time) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(replayDir)
	if err != nil {
		return nil, err
	}

	newReplays := make([]os.FileInfo, 0)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if file.ModTime().After(since) {
			newReplays = append(newReplays, file)
		}
	}

	return newReplays, nil
}

func UploadReplayToS3(uploader *s3manager.Uploader, key string, file io.Reader) error {
	params := &s3manager.UploadInput{
		Bucket: aws.String(S3BucketName),
		Key:    &key,
		Body:   file,
	}

	_, err := uploader.Upload(params)
	return err
}

func GetUploadResult(key string) (string, error) {
	URL := UploadEndpoint + "&FileName=" + key
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

func UploadReplay(uploader *s3manager.Uploader, path string) (string, error) {
	UUID := uuid.NewV4()
	key := fmt.Sprintf("%s.StormReplay", UUID)

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	err = UploadReplayToS3(uploader, key, file)
	if err != nil {
		return "", err
	}

	result, err := GetUploadResult(key)
	return result, err
}

func UploadReplays(uploader *s3manager.Uploader, paths []string) ([]string, []error) {
	// todo: parallelise
	results := make([]string, len(paths))
	errors := make([]error, len(paths))
	for i, path := range paths {
		result, err := UploadReplay(uploader, path)
		if err != nil {
			errors[i] = err
		} else {
			results[i] = result
		}
	}

	return results, errors
}

func UploadNewReplays(replayDir string, since time.Time) (map[string]int, error) {
	newReplays, err := ListNewReplays(replayDir, since)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0)
	for _, replay := range newReplays {
		p := path.Join(replayDir, replay.Name())
		paths = append(paths, p)
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: StaticCredentials,
		Region:      aws.String(S3BucketRegion),
	}))
	uploader := s3manager.NewUploader(sess)

	resultsMap := make(map[string]int)
	results, errors := UploadReplays(uploader, paths)
	for i, err := range errors {
		if err != nil {
			if _, exists := resultsMap[err.Error()]; exists {
				resultsMap[err.Error()]++
			} else {
				resultsMap[err.Error()] = 1
			}
		} else {
			result := results[i]
			if _, exists := resultsMap[result]; exists {
				resultsMap[result]++
			} else {
				resultsMap[result] = 1
			}
		}
	}

	return resultsMap, nil
}

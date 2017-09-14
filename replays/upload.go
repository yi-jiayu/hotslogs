package replays

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/mitchellh/go-homedir"
	"github.com/satori/go.uuid"
)

const (
	accessKeyID     = "AKIAIESBHEUH4KAAG4UA"
	secretAccessKey = "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"
	S3BucketName    = "heroesreplays"
	S3BucketRegion  = "us-west-2"

	uploadEndpoint = "https://www.hotslogs.com/UploadFile?Source="
)

const (
	windowsDefaultReplayLocationGlob = "Documents/Heroes of the Storm/Accounts/*/*-Hero-*/Replays/Multiplayer/*"
	osxDefaultReplayLocationGlob     = "Library/Application Support/Blizzard/Heroes of the Storm/Accounts/########/#-Hero-#-######/Replays/Multiplayer/*"
)

var (
	StaticCredentials = credentials.NewStaticCredentials(accessKeyID, secretAccessKey, "")
)

func listReplaysInDefaultLocation() ([]string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	var pattern string
	switch runtime.GOOS {
	case "windows":
		pattern = filepath.Join(home, windowsDefaultReplayLocationGlob)
	case "darwin":
		pattern = filepath.Join(home, osxDefaultReplayLocationGlob)
	default:
		return nil, errors.New(fmt.Sprintf("os not supported (%s)", runtime.GOOS))
	}

	return filepath.Glob(pattern)
}

func listNewReplaysInDefaultLocation(since time.Time) ([]string, error) {
	replays, err := listReplaysInDefaultLocation()
	if err != nil {
		return nil, err
	}

	newReplays := make([]string, 0)
	for _, replay := range replays {
		fi, err := os.Stat(replay)
		if err != nil {
			continue
		}

		if fi.IsDir() {
			continue
		}

		if fi.ModTime().After(since) {
			newReplays = append(newReplays, replay)
		}
	}

	return newReplays, nil
}

func ListNewReplays(replayDir string, since time.Time) ([]string, error) {
	if replayDir == "" {
		return listNewReplaysInDefaultLocation(since)
	} else {
		files, err := ioutil.ReadDir(replayDir)
		if err != nil {
			return nil, err
		}

		newReplays := make([]string, 0)
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			if file.ModTime().After(since) {
				newReplays = append(newReplays, filepath.Join(replayDir, file.Name()))
			}
		}

		return newReplays, nil
	}
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
	URL := uploadEndpoint + "&FileName=" + key
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
	errs := make([]error, len(paths))
	for i, path := range paths {
		result, err := UploadReplay(uploader, path)
		if err != nil {
			errs[i] = err
		} else {
			results[i] = result
		}
	}

	return results, errs
}

func UploadNewReplays(replayDir string, since time.Time) (map[string]int, error) {
	newReplays, err := ListNewReplays(replayDir, since)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0)
	for _, replay := range newReplays {
		path := filepath.Join(replayDir, replay)
		paths = append(paths, path)
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: StaticCredentials,
		Region:      aws.String(S3BucketRegion),
	}))
	uploader := s3manager.NewUploader(sess)

	resultsMap := make(map[string]int)
	results, errs := UploadReplays(uploader, paths)
	for i, err := range errs {
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

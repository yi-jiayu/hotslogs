package replays

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/mitchellh/go-homedir"
)

const (
	windowsDefaultReplayLocationGlob = "Documents/Heroes of the Storm/Accounts/*/*-Hero-*/Replays/Multiplayer/*.StormReplay"
	osxDefaultReplayLocationGlob     = "Library/Application Support/Blizzard/Heroes of the Storm/Accounts/########/#-Hero-#-######/Replays/Multiplayer/*.StormReplay"
)

func list(home string) ([]string, error) {
	var pattern string
	switch runtime.GOOS {
	case "windows":
		pattern = filepath.Join(home, windowsDefaultReplayLocationGlob)
	case "darwin":
		pattern = filepath.Join(home, osxDefaultReplayLocationGlob)
	default:
		return nil, errors.New(fmt.Sprintf("os not supported (%s)", runtime.GOOS))
	}

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	replays := make([]string, 0)
	for _, match := range matches {
		fi, err := os.Stat(match)
		if err != nil {
			continue
		}

		if !fi.IsDir() {
			replays = append(replays, match)
		}
	}

	return replays, nil
}

func List() ([]string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	return list(home)
}

func Since(t time.Time) ([]string, error) {
	replays, err := List()
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

		if fi.ModTime().After(t) {
			newReplays = append(newReplays, replay)
		}
	}

	return newReplays, nil
}

func ListNewReplays(replayDir string, since time.Time) ([]string, error) {
	if replayDir == "" {
		return Since(since)
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

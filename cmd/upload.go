// Copyright © 2017 Jiayu Yi <yi-jiayu@users.noreply.github.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yi-jiayu/hotslogs/hotsapi"
	"github.com/yi-jiayu/hotslogs/hotslogs"
	"github.com/yi-jiayu/hotslogs/replays"
)

const TimeFormat = "15:04"

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:     "upload",
	Aliases: []string{"up"},
	Short:   "Upload new replays",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// check if $HOME/.hotslogs folder exists
		home, err := homedir.Dir()
		if err != nil {
			panic(err)
		}

		datadir := filepath.Join(home, ".hotslogs")
		fi, err := os.Stat(datadir)
		if err != nil {
			if !os.IsNotExist(err) {
				panic(err)
			}

			os.Mkdir(datadir, 0600)
		} else {
			if !fi.IsDir() {
				panic("Oh no, $HOME/.hotslogs already exists but is not a directory.")
			}
		}

		db, err := bolt.Open(filepath.Join(datadir, "data"), 0600, nil)
		if err != nil {
			panic(err)
		}

		destinations := viper.GetStringSlice("destinations")
		replays_, err := replays.List()
		for _, dest := range destinations {
			switch dest {
			case "hotslogs":
				UploadToHOTSLogs(db, replays_)
			case "hotsapi":
				UploadToHotSAPI(db, replays_)
			default:
				fmt.Printf("Oh no, invalid upload destination: %s. Skipping...", dest)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(uploadCmd)
}

func FilterUploadedReplays(db *bolt.DB, bucket string, replays []string) []string {
	newReplays := make([]string, 0)
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			newReplays = replays
			return nil
		}

		for _, replay := range replays {
			if b.Get([]byte(replay)) == nil {
				newReplays = append(newReplays, replay)
			}
		}

		return nil
	})

	return newReplays
}

func UploadToHOTSLogs(db *bolt.DB, replays []string) {
	fmt.Printf("[%s] Starting upload to HOTS Logs\n", time.Now().Format(TimeFormat))

	newReplays := FilterUploadedReplays(db, "hotslogs", replays)
	if len(newReplays) == 0 {
		fmt.Println("No new replays, nothing to upload.")
		fmt.Printf("[%s] Finished upload to HOTS Logs\n", time.Now().Format(TimeFormat))
		return
	} else {
		fmt.Printf("Found %d new replays.\n", len(newReplays))
	}

	uploader := hotslogs.NewUploader()

	for i, replay := range replays {
		err := db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("hotslogs"))
			if err != nil {
				return err
			}

			fmt.Printf("(%d of %d) Uploading %s...", i, len(newReplays), filepath.Base(replay))

			result, err := uploader.UploadReplay(replay)
			if err != nil {
				fmt.Printf(`\rerror uploading "%s": %v\n`, replay, err)
				return nil
			}

			fmt.Printf("\rUploaded %s (%s)\n", filepath.Base(replay), result)

			err = b.Put([]byte(replay), []byte(result))
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			fmt.Printf("database error: %v\n", err)
		}
	}

	fmt.Printf("[%s] Finished upload to HOTS Logs\n", time.Now().Format(TimeFormat))
}

func UploadToHotSAPI(db *bolt.DB, replays []string) {
	fmt.Printf("[%s] Starting upload to HotS API\n", time.Now().Format(TimeFormat))

	newReplays := FilterUploadedReplays(db, "hotsapi", replays)
	if len(newReplays) == 0 {
		fmt.Println("No new replays, nothing to upload.")
		fmt.Printf("[%s] Finished upload to HotS API\n", time.Now().Format(TimeFormat))
		return
	} else {
		fmt.Printf("Found %d new replays.\n", len(newReplays))
	}

	uploader := hotsapi.NewUploader()

	for i, replay := range newReplays {
		err := db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("hotsapi"))
			if err != nil {
				return err
			}

			fmt.Printf("(%d of %d) Uploading %s...", i, len(newReplays), filepath.Base(replay))

			result, err := uploader.UploadReplay(replay)
			if err != nil {
				fmt.Printf(`\rerror uploading "%s": %v\n`, replay, err)
				return nil
			}

			fmt.Printf("\rUploaded %s (%s)\n", filepath.Base(replay), result.Status)

			err = b.Put([]byte(replay), []byte(result.Status))
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			fmt.Printf("database error: %v\n", err)
		}
	}

	fmt.Printf("[%s] Finished upload to HotS API\n", time.Now().Format(TimeFormat))
}

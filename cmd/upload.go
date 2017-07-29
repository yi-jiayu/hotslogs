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
	"log"
	"os"
	"path"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yi-jiayu/hotslogs/hotslogs"
)

var (
	dryRun bool
)

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
		replayDir := viper.GetString("replayDir")
		if replayDir == "" {
			log.Fatal("err: replay dir not set")
		}

		lastUploadTime := viper.GetTime("lastUploadTime")
		fmt.Printf("Looking for new replays since: %s\n", lastUploadTime)

		now := time.Now()
		newReplays, err := hotslogs.ListNewReplays(replayDir, lastUploadTime)
		if err != nil {
			log.Fatal(err)
		}

		if len(newReplays) == 0 {
			fmt.Println("No new replays since last upload, exiting.")
			os.Exit(0)
		} else {
			fmt.Printf("Found %d new replays since last upload.\n", len(newReplays))
		}

		fmt.Println("Uploading new replays...")

		sess := session.Must(session.NewSession(&aws.Config{
			Credentials: hotslogs.StaticCredentials,
			Region:      aws.String(hotslogs.S3BucketRegion),
		}))
		uploader := s3manager.NewUploader(sess)

		for _, replay := range newReplays {
			fmt.Printf("  %s: ", replay.Name())
			if !dryRun {
				result, err := hotslogs.UploadReplay(uploader, path.Join(replayDir, replay.Name()))
				if err != nil {
					fmt.Printf("ERROR (%s)\n", err)
				} else {
					fmt.Printf("DONE (%s)\n", result)
				}
			} else {
				fmt.Println("SKIPPED (Dry run)")
			}
		}

		if !dryRun {
			// update config file
			fmt.Print("Updating config file... ")

			config := fmt.Sprintf("replayDir: %s\nlastUploadTime: %s\n", replayDir, now.Format(time.RFC3339))
			file, err := os.OpenFile(viper.ConfigFileUsed(), os.O_TRUNC|os.O_CREATE, 0666)
			if err != nil {
				log.Fatal(err)
			}

			_, err = file.WriteString(config)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Done.")
		}
	},
}

func init() {
	RootCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	uploadCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Dry run")
}

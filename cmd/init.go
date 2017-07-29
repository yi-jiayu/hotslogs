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
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise hotslog config",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		configFile := viper.ConfigFileUsed()

		if configFile == "" {
			home, err := homedir.Dir()
			if err != nil {
				panic(err)
			}
			configFile = filepath.Join(home, ".hotslogs.yaml")
		}

		if fi, err := os.Stat(configFile); !os.IsNotExist(err) {
			if fi.IsDir() {
				fmt.Fprintf(os.Stderr, "error: %s is a directory.", configFile)
				os.Exit(1)
			}

			fmt.Printf("warning: %s already exists and and will be overwritten. Continue? (y/N) ", viper.ConfigFileUsed())
			input, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}

			input = strings.TrimSpace(input)
			if input != "y" && input != "Y" {
				fmt.Println("init aborted.")
				os.Exit(0)
			}
		}

		replayDirGuess := GuessReplayDir()
		if len(replayDirGuess) > 0 {
			fmt.Printf("Replay directory: (%s) ", replayDirGuess[0])
		} else {
			fmt.Printf("Replay directory: ")
		}

		var replayDir string
		input, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		input = strings.TrimSpace(input)

		if input == "" && replayDirGuess[0] != "" {
			replayDir = replayDirGuess[0]
		} else {
			replayDir = input
		}

		if fi, err := os.Stat(replayDir); !os.IsNotExist(err) {
			if !fi.IsDir() {
				fmt.Fprintf(os.Stderr, "error: %s exists but is not a directory.", replayDir)
				os.Exit(1)
			}
		} else {
			fmt.Fprintf(os.Stderr, "error: %s does not exist.", replayDir)
			os.Exit(1)
		}

		fmt.Printf("About to write to %s:\n", configFile)

		config := fmt.Sprintf("replayDir: %s\n", replayDir)
		fmt.Printf("\"\n%s\"\n", config)

		fmt.Print("Is this ok? (Y/n) ")
		input, err = reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		input = strings.TrimSpace(input)
		if input != "y" && input != "Y" && input != "" {
			fmt.Println("init aborted.")
			os.Exit(0)
		}

		fmt.Printf("Initialising %s... ", configFile)

		file, err := os.OpenFile(configFile, os.O_TRUNC|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}

		_, err = file.WriteString(config)
		if err != nil {
			panic(err)
		}

		fmt.Println("Done.")
	},
}

func GuessReplayDir() []string {
	home, err := homedir.Dir()
	if err != nil {
		return nil
	}

	switch runtime.GOOS {
	case "windows":
		glob := filepath.Join(home, "Documents/Heroes of the Storm/Accounts/*/*-Hero-*/Replays/Multiplayer")
		matches, err := filepath.Glob(glob)
		if err != nil {
			return nil
		}

		return matches
	case "darwin":
		// todo: macos support
		return nil
	default:
		return nil
	}
}

func init() {
	configCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var lcdir string

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//clone the remote repository to local place
		gitUrl := "https://github.com/" + username + "/" + repo + ".git"
		fmt.Println(gitUrl)

		cloneCmd := exec.Command("git", "clone", gitUrl)
		err := cloneCmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		//copy out all file in repo under local repository
		absLocalPath, err := filepath.Abs(lcdir)
		if err != nil {
			log.Fatal(err)
		}

		err = os.Chdir(repo)
		if err != nil {
			log.Fatal(err)
		}

		mvCmd := exec.Command("/bin/sh", "-c", "mv * "+absLocalPath)
		err = mvCmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		//push the empty repo to remote
		//do git push in local repository
		addCmd := exec.Command("git", "add", "*")
		err = addCmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
		index := rnd.Intn(10000000)
		indexstr := strconv.Itoa(index)
		cmtstr := "commit" + indexstr

		cmtCmd := exec.Command("git", "commit", "-m", cmtstr)
		err = cmtCmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		//delete local repository
		err = os.Chdir("..")
		if err != nil {
			log.Fatal(err)
		}

		err = os.RemoveAll(repo)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Download file finish!")
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.
	downloadCmd.PersistentFlags().StringVarP(&lcdir, "lcdir", "l", ".", "the local directory you want to store the downloaded file")
	downloadCmd.PersistentFlags().StringVarP(&repo, "repo", "r", "thunderfile", "the name of remote github repo")
	downloadCmd.PersistentFlags().StringVarP(&username, "uname", "u", "hjiangsse", "the github user name")
}

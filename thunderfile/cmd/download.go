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
	"strings"
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
		gitSsh := "git@github.com:" + username + "/" + repo + ".git"
		fmt.Println("[Repo: " + gitSsh + "]")
		
		cloneCmd := exec.Command("git", "clone", gitSsh)
		err := cloneCmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("[download: repo clone finished!]")
		
		//copy out all file in repo under local repository
		absLocalPath, err := filepath.Abs(lcdir)
		if err != nil {
			log.Fatal(err)
		}
	
		err = os.Chdir(repo)
		if err != nil {
			log.Fatal(err)
		}

		//walk through all files under current directory, move it
		//to the dest place
		err = filepath.Walk(".",
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					log.Fatal(err)
					return err
				}

				if !strings.Contains(path, ".git") && path[0] != '.' {
					destpath := filepath.Join(absLocalPath, path)
					os.Rename(path, destpath)
					os.Remove(path)
				}
				return nil
			})
			
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("[download: copy repo file to local finish!]")
		
		//push the empty repo to remote
		//do git push in local repository
		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
		index := rnd.Intn(10000000)
		indexstr := strconv.Itoa(index)
		cmtstr := "commit" + indexstr

		cmtCmd := exec.Command("git", "commit", "-a", "-m", cmtstr)
		err = cmtCmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		
		pushCmd := exec.Command("git", "push")
		err = pushCmd.Run()
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

		fmt.Println("[Download: Download file finish! Update remote repository OK!]")
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.
	downloadCmd.PersistentFlags().StringVarP(&lcdir, "lcdir", "l", ".", "the local directory you want to store the downloaded file")
	downloadCmd.PersistentFlags().StringVarP(&repo, "repo", "r", "thunderfile", "the name of remote github repo")
	downloadCmd.PersistentFlags().StringVarP(&username, "uname", "u", "hjiangsse", "the github user name")
}

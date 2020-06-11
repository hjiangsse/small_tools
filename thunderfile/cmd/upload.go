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
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/go-github/v32/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var repo string
var lcfile string
var username string

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload your file to a remote reposity",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		token := os.Getenv("GITHUB_AUTH_TOKEN")
		if token == "" {
			fmt.Println("Please set your github auth token as environment variable *GITHUB_AUTH_TOKEN*")
			os.Exit(1)
		}

		//initialize github client
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		tc := oauth2.NewClient(ctx, ts)
		client := github.NewClient(tc)

		//if repo not exist, create it
		private := false
		desc := "sync files just like a thunder"
		r := &github.Repository{Name: &repo, Private: &private, Description: &desc}
		rep, _, err := client.Repositories.Create(ctx, "", r)
		if err != nil {
			fmt.Println("Create Repo Failed!--Alread exist.")
		} else {
			fmt.Printf("Successfully created new repo: %v\n", rep.GetName())
		}

		//clone the repo to local place
		gitSsh := "git@github.com:" + username + "/" + repo + ".git"
		fmt.Println("[Upload: " + gitSsh + "]")

		cloneCmd := exec.Command("git", "clone", gitSsh)
		err = cloneCmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		//copy all the file(s) you want to upload to local repo
		absLocalPath, err := filepath.Abs(lcfile)
		if err != nil {
			log.Fatal(err)
		}
		
		var cpCmd *exec.Cmd
		pathInfo, err := os.Stat(absLocalPath)
		if err != nil {
			log.Fatal(err)
		}

		if pathInfo.IsDir() {
			cpCmd = exec.Command("cp", "-r", absLocalPath, repo)
		} else {
			cpCmd = exec.Command("cp", absLocalPath, repo)
		}

		err = cpCmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		//do git push in local repository
		err = os.Chdir(repo)
		if err != nil {
			log.Fatal(err)
		}

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

		fmt.Println("Upload file finish!")
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	uploadCmd.PersistentFlags().StringVarP(&repo, "repo", "r", "thunderfile", "the name of remote github repo")
	uploadCmd.PersistentFlags().StringVarP(&username, "uname", "u", "hjiangsse", "the github user name")
	uploadCmd.PersistentFlags().StringVarP(&lcfile, "lcfile", "l", "./README.md", "the local file you want to upload")
}

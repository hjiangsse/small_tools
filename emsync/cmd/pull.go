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
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

const (
	LOCALREPOSTR string = "localrepo"
)

var (
	localrepo string
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "pull the latest github emacs configs",
	Long:  `pull the latest emacs config files from github`,
	Run: func(cmd *cobra.Command, args []string) {
		//change to local repo directory
		if _, err := os.Stat(path.Clean(localrepo)); os.IsNotExist(err) {
			fmt.Printf("%v is not exist! Bye!\n", localrepo)
			return
		}

		if err := os.Chdir(localrepo); err != nil {
			fmt.Println("Change working dir error: ", err.Error())
			return
		}

		//do clean and pull
		cleanCmd := exec.Command("git", "clean", "-f")
		fmt.Println("Running git clean -f...")
		err := cleanCmd.Run()
		if err != nil {
			fmt.Println("Git clean error: ", err.Error())
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	pullCmd.PersistentFlags().StringVarP(&localrepo, LOCALREPOSTR, "l", "~/github/emacsconf/", "local emacs config repo path")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

//given a user command inout path, change it into os absolute path
func GetAbsLocalCnfPath(inputPath string) (string, error) {
	if strings.Contains(inputPath, "~") {

	}
}

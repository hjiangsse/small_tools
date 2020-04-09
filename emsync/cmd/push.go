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

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "push local emacs configs to github",
	Long:  `push local emacs configs to github`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		//copy local emacs configs to local repo
		homedir, err := os.UserHomeDir()
		if err != nil {
			return
		}

		absRepo, err := GetAbsLocalCnfPath(localrepo)
		if err != nil {
			return
		}

		localPath := path.Join(homedir, ".emacs.d")
		err = DoCopyFromLocalToRepo(localPath, absRepo)
		if err != nil {
			return
		}

		//change current work directory to local github repo
		if err := os.Chdir(absRepo); err != nil {
			return
		}
		fmt.Printf("change working dir to : %v\n", absRepo)

		//do git command(add *, commit, push)
		gitAddCmd := exec.Command("git", "add", "*")
		if err := gitAddCmd.Run(); err != nil {
			return
		}
		fmt.Println("git add * finish!")

		/*
			var cmdout bytes.Buffer
			var cmderr bytes.Buffer

			gitCommitCmd.Stdout = &cmdout
			gitCommitCmd.Stderr = &cmderr
			if err := gitCommitCmd.Run(); err != nil {
				fmt.Printf(fmt.Sprint(err) + ": " + cmderr.String())
				return
			}
		*/

		gitCommitCmd := exec.Command("git", "commit", "-m", "commit")
		res, err := gitCommitCmd.Output()
		if err != nil && !strings.Contains(string(res), "up to date") {
			fmt.Println(err)
			return
		}
		fmt.Println("git commit finish!")

		gitPushCmd := exec.Command("git", "push")
		if err := gitPushCmd.Run(); err != nil {
			return
		}
		fmt.Println("git push finish!")
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
	//pushCmd.PersistentFlags().StringVarP(&localrepo, LOCALREPOSTR, "l", "~/github/emacsconf/", "local emacs config repo path")
}

//copy local emacs configs to local github repository
func DoCopyFromLocalToRepo(localPath string, destPath string) error {
	var err error

	localInitPath := path.Join(localPath, "init.el")
	destInitPath := path.Join(destPath, "inits")

	localModulesPath := path.Join(localPath, "modules")
	destModulesPath := path.Join(destPath, "modules")

	localTempsPath := path.Join(localPath, "templates")
	destTempsPath := path.Join(destPath, "templates")

	localUtilsPath := path.Join(localPath, "utils")
	destUtilsPath := path.Join(destPath, "utils")

	//copy emacs init file
	cpInitCmd := exec.Command("cp", localInitPath, destInitPath)
	err = cpInitCmd.Run()
	if err != nil {
		return err
	}
	fmt.Printf("copy %v to %v finish!\n", localInitPath, destInitPath)

	//copy modules files to local github repo
	err = copyTo(localModulesPath, destModulesPath, ".el")
	if err != nil {
		fmt.Println("copy modules files error!")
		return err
	}

	//copy templates files to local github repo
	err = copyTo(localTempsPath, destTempsPath, ".txt")
	if err != nil {
		fmt.Println("copy templates files error!")
		return err
	}

	//copy utils files to local github repo
	err = copyTo(localUtilsPath, destUtilsPath, ".el")
	if err != nil {
		fmt.Println("copy utils files error!")
		return err
	}

	return nil
}

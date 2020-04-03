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
	"io/ioutil"
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
		//preprocess local repo path
		clrRepo, err := GetAbsLocalCnfPath(path.Clean(localrepo))
		if err != nil {
			fmt.Printf("Preprocess local repository path error! Bye!")
			return
		}

		//change to local repo directory
		if _, err := os.Stat(clrRepo); os.IsNotExist(err) {
			fmt.Printf("%v is not exist! Bye!\n", clrRepo)
			return
		}

		if err := os.Chdir(clrRepo); err != nil {
			fmt.Println("Change working dir error: ", err.Error())
			return
		}

		//check if whether exists .git file in this directory
		lsCmd := exec.Command("ls", "-a")
		lsOut, err := lsCmd.Output()
		if err != nil {
			fmt.Println("ls -a error!")
			return
		}

		if !strings.Contains(string(lsOut), ".git") {
			fmt.Println("local repo dir has no .git file, error!")
			return
		}

		//do clean and pull
		cleanCmd := exec.Command("git", "clean", "-f")
		fmt.Println("Running git clean -f...")
		err = cleanCmd.Run()
		if err != nil {
			fmt.Println("Git clean error: ", err.Error())
			return
		}

		pullCmd := exec.Command("git", "pull")
		fmt.Println("Running git pull...")
		err = pullCmd.Run()
		if err != nil {
			fmt.Println("Git pull error: ", err.Error())
			return
		}

		//copy the latest file to local .emacs.d directory
		err = DoCopy(".emacs.d")
		if err != nil {
			fmt.Println("copy file to .emacs.d error: ", err.Error())
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
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return strings.Replace(inputPath, "~", home, 1), nil
	}

	return inputPath, nil
}

//copy file to $HOME/destPath
func DoCopy(destPath string) error {
	var err error
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	destDirPath := path.Join(homedir, destPath)
	initFilePath := path.Join(cwd, "inits/init.el")

	modulesFilePath := path.Join(cwd, "modules")
	modulesFileDstPath := path.Join(destDirPath, "modules/")

	tempFilePath := path.Join(cwd, "templates")
	tempFileDstPath := path.Join(destDirPath, "templates/")

	utilsFilePath := path.Join(cwd, "utils")
	utilsFileDstPath := path.Join(destDirPath, "utils/")

	//copy emacs init file
	cpInitCmd := exec.Command("cp", initFilePath, destDirPath)
	err = cpInitCmd.Run()
	if err != nil {
		return err
	}

	//copy modules file
	err = copyTo(modulesFilePath, modulesFileDstPath, ".el")
	if err != nil {
		return err
	}

	//copy template file
	err = copyTo(tempFilePath, tempFileDstPath, ".txt")
	if err != nil {
		return err
	}

	//copy utils file
	err = copyTo(utilsFilePath, utilsFileDstPath, ".el")
	if err != nil {
		return err
	}

	return nil
}

//copy all file in "source", which have "suffix" in the tail,
//to "dest", example: copyTo("../test/source", "../test/dest", ".txt")
func copyTo(source, dest, suffix string) error {
	files, err := ioutil.ReadDir(source)
	if err != nil {
		return err
	}

	for _, f := range files {

		if strings.Contains(f.Name(), suffix) {
			cpCmd := exec.Command("cp", path.Join(source, f.Name()), dest)
			err := cpCmd.Run()
			if err != nil {
				return err
			}

			fmt.Printf("cp %v to %v\n", path.Join(source, f.Name()), dest)
		}
	}

	return nil
}

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
	"regexp"

	"github.com/spf13/cobra"
)

var (
	inputpath  string
	outputpath string
)

var (
	titileLine = regexp.MustCompile(`^(\**) .*$`)
	sourceLine = regexp.MustCompile(`\s*#\+(BEGIN|END)_SRC.*`)
)

const (
	S_BEGIN_SPACES = iota
	S_SOURCE_START
	S_BEGIN_SOURCE
	S_BEGIN_SOURCE_END
	S_END_SOURCE
	S_END_SOURCE_END
)

// tranCmd represents the tran command
var tranCmd = &cobra.Command{
	Use:   "tran",
	Short: "org file transform to md file",
	Long: `jekyll use markdown file as post page, so i need a
tool to transform org file to markdown file`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("input path: ", inputpath)
		fmt.Println("output path: ", outputpath)

		testLine := []byte("**** This is a test line")
		if isOrgTitle(testLine) {
			fmt.Println(string(transTitleLine(testLine)))
		}

		testSource1 := []byte("   #+BEGIN_SRC go  ")
		if isOrgSouce(testSource1) {
			fmt.Println(string(transSourceLine(testSource1)))
		}

		testSource2 := []byte("#+BEGIN_SRC go  ")
		if isOrgSouce(testSource2) {
			fmt.Println(string(transSourceLine(testSource2)))
		}

		testSource3 := []byte("#+END_SRC go  ")
		if isOrgSouce(testSource3) {
			fmt.Println(string(transSourceLine(testSource3)))
		}
	},
}

func init() {
	rootCmd.AddCommand(tranCmd)

	// Here you will define your flags and configuration settings.
	tranCmd.PersistentFlags().StringVarP(&inputpath, "input-file", "i", "input.org", "the path of the file for transform")
	tranCmd.PersistentFlags().StringVarP(&outputpath, "output-file", "o", "output.md", "the path of the file for output")
}

// tell if a line is org title
func isOrgTitle(line []byte) bool {
	if titileLine.Match(line) {
		return true
	}
	return false
}

// tell if a line is org source code snippet
func isOrgSouce(line []byte) bool {
	if sourceLine.Match(line) {
		return true
	}
	return false
}

// transform a org title to markdown version
func transTitleLine(line []byte) []byte {
	for i, e := range line {
		if e == ' ' {
			return line
		}

		if e == '*' {
			line[i] = '#'
		}
	}
	return line
}

// transform a org source snippt line to markdown version
func transSourceLine(line []byte) []byte {
	res := make([]byte, 0)
	state := S_BEGIN_SPACES

	for _, e := range line {
		switch state {
		case S_BEGIN_SPACES:
			if e == ' ' {
				res = append(res, e)
			} else {
				state = S_SOURCE_START
			}
		case S_SOURCE_START:
			if e == 'B' {
				state = S_BEGIN_SOURCE
			}

			if e == 'E' {
				state = S_END_SOURCE
			}
		case S_BEGIN_SOURCE:
			if e == ' ' {
				res = append(res, []byte("``` ")...)
				state = S_BEGIN_SOURCE_END
			}
		case S_END_SOURCE:
			if e == ' ' {
				res = append(res, []byte("``` ")...)
				state = S_END_SOURCE_END
			}
		case S_BEGIN_SOURCE_END:
			res = append(res, e)
		case S_END_SOURCE_END:
			res = append(res, e)
		default:
			res = append(res, e)
		}
	}

	return res
}

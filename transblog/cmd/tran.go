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
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var (
	inputpath  string
	outputpath string
	srcimgpath string
	dstimgpath string
	blogpath   string
)

var (
	titileLine = regexp.MustCompile(`^(\*+) .*$`)
	sourceLine = regexp.MustCompile(`\s*#\+(BEGIN|END)_SRC.*`)
	imageLine  = regexp.MustCompile(`\s*\[\[file:(.*)\]\[(.*)\]\]\s*`)
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
	Long: `jekyll use markdown file as post page, so i need a tool to
transform org file to markdown file
           usage:
           transblog tran -i input.org -o output.md -d image_url_path`,
	Run: func(cmd *cobra.Command, args []string) {
		//open the source org file
		infile, err := os.Open(inputpath)
		if err != nil {
			log.Printf("open input file error: %v\n", err)
			os.Exit(1)
		}
		defer infile.Close()

		//create or open the dest md file
		outfile, err := os.Create(outputpath)
		if err != nil {
			log.Printf("open output file error: %v\n", err)
			os.Exit(1)
		}
		defer outfile.Close()

		//trans source file line by line into dest file
		err = transformFile(infile, outfile)
		if err != nil {
			log.Printf("transform %v to %v failed!\n", infile.Name(), outfile.Name())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(tranCmd)

	// Here you will define your flags and configuration settings.
	tranCmd.PersistentFlags().StringVarP(&inputpath, "input-file", "i", "input.org", "the path of the file for transform")
	tranCmd.PersistentFlags().StringVarP(&outputpath, "output-file", "o", "output.md", "the path of the file for output")
	tranCmd.PersistentFlags().StringVarP(&srcimgpath, "srcimg-path", "s", "/Users/hjiang/github/orgnization/graph", "the source image file")
	tranCmd.PersistentFlags().StringVarP(&dstimgpath, "dstimg-path", "d", "/assets/gc", "the dest image path")
	tranCmd.PersistentFlags().StringVarP(&blogpath, "blog-path", "b", "/Users/hjiang/github/myblog/hjiangsse.github.io", "the blog path")
}

//transform source (.org) file into dest (.md) file
func transformFile(src *os.File, dst *os.File) error {
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		curLineSlice := []byte(scanner.Text())

		//if current line is org title line
		if isOrgTitle(curLineSlice) {
			fmt.Println(string(transTitleLine(curLineSlice)))
			dst.WriteString(string(transTitleLine(curLineSlice)) + "\n")
			continue
		}

		//if current line is code snippet begin or end
		if isOrgSouce(curLineSlice) {
			fmt.Println(string(transSourceLine(curLineSlice)))
			dst.WriteString(string(transSourceLine(curLineSlice)) + "\n")
			continue
		}

		//if current line is image file line
		if isInsImage(curLineSlice) {
			mdInsLineSlice, err := transImageLine(curLineSlice)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
			fmt.Println(string(mdInsLineSlice))
			dst.WriteString(string(mdInsLineSlice) + "\n")
			continue
		}

		//other normal lines
		dst.WriteString(string(curLineSlice) + "\n")
	}
	return nil
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

// tell if a line stands org inserted image
// example: [[file:graph/ack_send.png][send five tasks to rebbitMq]]
func isInsImage(line []byte) bool {
	if imageLine.Match(line) {
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
			if e == 'C' {
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

// transform org image line to markdown version
func transImageLine(line []byte) ([]byte, error) {
	trimedLine := strings.Trim(string(line), "[] ")
	trimedSegs := strings.Split(trimedLine, "][")

	if len(trimedSegs) > 2 || len(trimedSegs) == 0 {
		return nil, errors.New("invalid org image line, two many segments")
	}

	fileSegs := strings.Split(trimedSegs[0], ":")
	if len(fileSegs) != 2 {
		return nil, errors.New("invalid org file path")
	}

	fileName := filepath.Base(fileSegs[1])
	mdFileName := filepath.Join(dstimgpath, fileName)

	srcimgurl := filepath.Join(srcimgpath, fileName)
	dstimgpath := filepath.Join(blogpath, dstimgpath)

	if _, err := os.Stat(dstimgpath); os.IsNotExist(err) {
		mkcmd := exec.Command("mkdir", dstimgpath)
		err := mkcmd.Run()
		if err != nil {
			return []byte(""), err
		}
	}

	cpcmd := exec.Command("cp", srcimgurl, dstimgpath)
	err := cpcmd.Run()
	if err != nil {
		return []byte(""), err
	}

	return []byte("![" + trimedSegs[1] + "](" + mdFileName + ")"), nil
}

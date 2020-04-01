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
	"path"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

const (
	BACKDIRSTR string = "backdir"
	OUTDIRSTR  string = "outdir"
	TOPICSTR   string = "topic"
)

var (
	backdir string
	topic   string
	outdir  string
)

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "dump nsqd backup data into json file",
	Long:  `dump nsqd backup data into json file`,
	Run: func(cmd *cobra.Command, args []string) {
		//if user want to dump data of a specific topic
		if cmd.Flags().Changed(TOPICSTR) {
			err := DumpSpecificTopic(backdir, outdir, topic)
			if err != nil {
				panic(err)
			}
		} else {
			var wg sync.WaitGroup

			//dump all topic data
			topics, err := GetAllTopics(backdir)
			if err != nil {
				panic(err)
			}

			for _, topic := range topics {
				wg.Add(1)
				go func(tp string) {
					defer wg.Done()
					DumpSpecificTopic(backdir, outdir, tp)
				}(topic)
			}

			wg.Wait()
		}
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)

	dumpCmd.PersistentFlags().StringVarP(&backdir, BACKDIRSTR, "b", ".", "nsq back file dir path")
	dumpCmd.PersistentFlags().StringVarP(&outdir, OUTDIRSTR, "o", ".", "nsq dumped file dir path")
	dumpCmd.PersistentFlags().StringVarP(&topic, TOPICSTR, "t", "test", "the specific topic you want to dump")
}

//dump specific topic data under a path into json file to outpath
func DumpSpecificTopic(datapath, outpath, topic string) error {
	var filenum uint64
	var err error
	var topicfile string
	var dumpfile string

	cleanedDataPath := path.Clean(datapath)
	cleanedOutPath := path.Clean(outpath)

	topicfile = GetDataFileName(cleanedDataPath, topic, filenum)
	dumpfile = GetDumpFileName(cleanedOutPath, topic, filenum)

	//if topic file not exist
	_, err = os.Stat(topicfile)
	if os.IsNotExist(err) {
		return err
	}

	//if outpath not exist, create it
	if _, err = os.Stat(cleanedOutPath); os.IsNotExist(err) {
		err = os.Mkdir(cleanedOutPath, 0766)
		if err != nil {
			return err
		}
	}

	for {
		//parse and dump each file belong to current topic
		err = DecodeAndDumpNsqBackFile(topicfile, dumpfile)
		if err != nil {
			break
		}

		filenum++

		topicfile = GetDataFileName(cleanedDataPath, topic, filenum)
		dumpfile = GetDumpFileName(cleanedOutPath, topic, filenum)

		_, err = os.Stat(topicfile)
		if os.IsNotExist(err) {
			break
		}
	}

	return nil
}

//given path, topic and filenum, return nsqd .dat file name
func GetDataFileName(datapath, topic string, filenum uint64) string {
	return fmt.Sprintf(path.Join(datapath, "%s.diskqueue.%06d.dat"), topic, filenum)
}

//given path, topic and filenum, return dumped .json file name
func GetDumpFileName(outpath, topic string, filenum uint64) string {
	return fmt.Sprintf(path.Join(outpath, "%s.diskqueue.%06d.json"), topic, filenum)
}

//get all topic name under a specific directory
func GetAllTopics(datapath string) ([]string, error) {
	var res []string
	topicMap := make(map[string]bool)

	files, err := ioutil.ReadDir(path.Clean(datapath))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		//this is a nsqd backup file
		if strings.Contains(file.Name(), ".diskqueue.") {
			namesegs := strings.Split(file.Name(), ".")
			topicMap[namesegs[0]] = true
		}
	}

	for key, _ := range topicMap {
		res = append(res, key)
	}

	return res, nil
}

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
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var metafile string
var outfile string
var recordnum int

type FieldInfo struct {
	FieldLen  int
	FieldName string
	FieldType string
	InitVal   string
}

type CfgInfos struct {
	FieldInfos []FieldInfo
}

// encodeCmd represents the encode command
var encodeCmd = &cobra.Command{
	Use:   "encode",
	Short: "generate a binary file according to meta data",
	Long: `According to the meta data provided by user, which contain
filedlen, filename, fieldtype and init value, generate a binary file`,
	Run: func(cmd *cobra.Command, args []string) {
		err := generateBinaryFile(metafile, outfile, recordnum)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(encodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encodeCmd.PersistentFlags().String("foo", "", "A help for foo")
	encodeCmd.PersistentFlags().StringVarP(&metafile, "meta", "m", "./configs/meta.json", "meta data file path")
	encodeCmd.PersistentFlags().StringVarP(&outfile, "out", "o", "./out.log", "output data file path")
	encodeCmd.PersistentFlags().IntVarP(&recordnum, "num", "n", 10, "number of records in the output file")
}

//generateBinaryFile: generate binary file according to the mata file
func generateBinaryFile(metadatapath, outbinarypath string, recordnum int) error {
	//decode the meta data
	var fieldinfos CfgInfos
	cnfJsonFile, err := os.Open(metadatapath)
	if err != nil {
		return err
	}
	cfgBytes, _ := ioutil.ReadAll(cnfJsonFile)

	err = json.Unmarshal(cfgBytes, &fieldinfos)
	if err != nil {
		return err
	}

	//open or create the text file for output
	outfile, err := os.Create(outbinarypath)
	if err != nil {
		return err
	}

	//generate records and write to bianry file
	var tempBuf bytes.Buffer
	var fmtStr string
	var fieldStr string

	for i := 0; i < recordnum; i++ {
		for _, e := range fieldinfos.FieldInfos {
			switch e.FieldType {
			case "int64":
				val, err := strconv.ParseInt(e.InitVal, 10, 64)
				if err != nil {
					return err
				}

				binary.Write(&tempBuf, binary.LittleEndian, &val)
			case "int32":
				val, err := strconv.ParseInt(e.InitVal, 10, 32)
				if err != nil {
					return err
				}

				trval := int32(val)
				binary.Write(&tempBuf, binary.LittleEndian, &trval)
			case "int16":
				val, err := strconv.ParseInt(e.InitVal, 10, 16)
				if err != nil {
					return err
				}

				trval := int16(val)
				binary.Write(&tempBuf, binary.LittleEndian, &trval)
			case "uint64":
				val, err := strconv.ParseUint(e.InitVal, 10, 64)
				if err != nil {
					return err
				}

				binary.Write(&tempBuf, binary.LittleEndian, &val)
			case "uint32":
				val, err := strconv.ParseUint(e.InitVal, 10, 32)
				if err != nil {
					return err
				}

				trval := uint32(val)
				binary.Write(&tempBuf, binary.LittleEndian, &trval)
			case "uint16":
				val, err := strconv.ParseUint(e.InitVal, 10, 16)
				if err != nil {
					return err
				}

				trval := uint16(val)
				binary.Write(&tempBuf, binary.LittleEndian, &trval)
			case "string":
				fmtStr = fmt.Sprintf("%%%ds", e.FieldLen)
				fieldStr = fmt.Sprintf(fmtStr, e.InitVal)
				binary.Write(&tempBuf, binary.LittleEndian, []byte(fieldStr))
			}
		}
		outfile.Write(tempBuf.Bytes())
		tempBuf.Reset()
	}
	return nil
}

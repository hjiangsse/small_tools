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
	"io"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

var binaryfile string

// decodeCmd represents the decode command
var decodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "Decode a binary file to text file",
	Long:  `Decode a bianry file into text file according to meta data`,
	Run: func(cmd *cobra.Command, args []string) {
		err := decodeBinaryFile(binaryfile, metafile, outfile)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(decodeCmd)

	decodeCmd.PersistentFlags().StringVarP(&metafile, "meta", "m", "./configs/meta.json", "meta data file path")
	decodeCmd.PersistentFlags().StringVarP(&binaryfile, "in", "i", "./in.log", "binary data file path")
	decodeCmd.PersistentFlags().StringVarP(&outfile, "out", "o", "./out.log", "output data file path")
}

//decodeBinaryFile: decode binary file to text file according to the field info
//in the json meta file
func decodeBinaryFile(binaryfilepath, metadatapath, outtextpath string) error {
	//open the binary file
	file, err := os.Open(binaryfilepath)
	if err != nil {
		return err
	}
	defer file.Close()

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
	outfile, err := os.Create(outtextpath)
	if err != nil {
		return err
	}

	//use a buffer to decode every line
	tempBuf := make([]byte, 1024)
	var decodeBuf bytes.Buffer //strore the decoded result
	var fmtStr string
	var fieldStr string
	for {
		for _, e := range fieldinfos.FieldInfos {
			n, err := file.Read(tempBuf[0:e.FieldLen])
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}

			switch e.FieldType {
			case "int64":
				var val int64
				binary.Read(bytes.NewBuffer(tempBuf[:n]), binary.LittleEndian, &val)
				fmtStr = fmt.Sprintf("%%%dd", e.FieldLen)
				fieldStr = fmt.Sprintf(fmtStr, val)
			case "int32":
				var val int32
				binary.Read(bytes.NewBuffer(tempBuf[:n]), binary.LittleEndian, &val)
				fmtStr = fmt.Sprintf("%%%dd", e.FieldLen)
				fieldStr = fmt.Sprintf(fmtStr, val)
			case "int16":
				var val int16
				binary.Read(bytes.NewBuffer(tempBuf[:n]), binary.LittleEndian, &val)
				fmtStr = fmt.Sprintf("%%%dd", e.FieldLen)
				fieldStr = fmt.Sprintf(fmtStr, val)
			case "uint64":
				var val uint64
				binary.Read(bytes.NewBuffer(tempBuf[:n]), binary.LittleEndian, &val)
				fmtStr = fmt.Sprintf("%%%dd", e.FieldLen)
				fieldStr = fmt.Sprintf(fmtStr, val)
			case "uint32":
				var val uint64
				binary.Read(bytes.NewBuffer(tempBuf[:n]), binary.LittleEndian, &val)
				fmtStr = fmt.Sprintf("%%%dd", e.FieldLen)
				fieldStr = fmt.Sprintf(fmtStr, val)
			case "uint16":
				var val uint16
				binary.Read(bytes.NewBuffer(tempBuf[:n]), binary.LittleEndian, &val)
				fmtStr = fmt.Sprintf("%%%dd", e.FieldLen)
				fieldStr = fmt.Sprintf(fmtStr, val)
			case "string":
				var val = make([]byte, e.FieldLen)
				binary.Read(bytes.NewBuffer(tempBuf[:n]), binary.LittleEndian, val)
				fmtStr = fmt.Sprintf("%%%ds", e.FieldLen)
				fieldStr = fmt.Sprintf(fmtStr, val)
			}

			decodeBuf.WriteString("|" + fieldStr)
		}
		decodeBuf.WriteTo(outfile)
		outfile.WriteString("\n")
	}
	return nil
}

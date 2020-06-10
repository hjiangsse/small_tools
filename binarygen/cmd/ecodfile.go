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
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var textfile string

// ecodfileCmd represents the ecodfile command
var ecodfileCmd = &cobra.Command{
	Use:   "ecodfile",
	Short: "Encode a text file into a binary file",
	Long:  `Encode a text file into a binary file, according to the meta data`,
	Run: func(cmd *cobra.Command, args []string) {
		err := generateBinaryFromText(metafile, textfile, binaryfile)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(ecodfileCmd)

	ecodfileCmd.PersistentFlags().StringVarP(&metafile, "meta", "m", "./configs/meta.json", "meta data file path")
	ecodfileCmd.PersistentFlags().StringVarP(&textfile, "in", "i", "./in.txt", "input text file path")
	ecodfileCmd.PersistentFlags().StringVarP(&binaryfile, "out", "o", "./out.log", "output binary file path")
}

//generateBinaryFromText: convert a text file to binary file
func generateBinaryFromText(metadatapath, textfilepath, outbinarypath string) error {
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

	//open or create the binary file for output
	outfile, err := os.Create(outbinarypath)
	if err != nil {
		return err
	}
	defer outfile.Close()

	//open the text file to get input
	infile, err := os.Open(textfilepath)
	if err != nil {
		return err
	}
	defer infile.Close()

	//generate records and write to bianry file
	var tempBuf bytes.Buffer
	var fmtStr string
	var fieldStr string

	scanner := bufio.NewScanner(infile)
	for scanner.Scan() {
		line := scanner.Text()
		validLine := line[1:]
		lineSegs := strings.Split(validLine, "|")
		for i, field := range lineSegs {
			fieldInfo := fieldinfos.FieldInfos[i]

			switch fieldInfo.FieldType {
			case "int64":
				val, err := strconv.ParseInt(strings.Trim(field, " "), 10, 64)
				if err != nil {
					return err
				}

				binary.Write(&tempBuf, binary.LittleEndian, &val)
			case "int32":
				val, err := strconv.ParseInt(strings.Trim(field, " "), 10, 32)
				if err != nil {
					return err
				}

				trval := int32(val)
				binary.Write(&tempBuf, binary.LittleEndian, &trval)
			case "int16":
				val, err := strconv.ParseInt(strings.Trim(field, " "), 10, 16)
				if err != nil {
					return err
				}

				trval := int16(val)
				binary.Write(&tempBuf, binary.LittleEndian, &trval)
			case "uint64":
				val, err := strconv.ParseUint(strings.Trim(field, " "), 10, 64)
				if err != nil {
					return err
				}

				binary.Write(&tempBuf, binary.LittleEndian, &val)
			case "uint32":
				val, err := strconv.ParseUint(strings.Trim(field, " "), 10, 32)
				if err != nil {
					return err
				}

				trval := uint32(val)
				binary.Write(&tempBuf, binary.LittleEndian, &trval)
			case "uint16":
				val, err := strconv.ParseUint(strings.Trim(field, " "), 10, 16)
				if err != nil {
					return err
				}

				trval := uint16(val)
				binary.Write(&tempBuf, binary.LittleEndian, &trval)
			case "string":
				fmtStr = fmt.Sprintf("%%%ds", fieldInfo.FieldLen)
				fieldStr = fmt.Sprintf(fmtStr, field)
				binary.Write(&tempBuf, binary.LittleEndian, []byte(fieldStr))
			}
		}
		outfile.Write(tempBuf.Bytes())
		tempBuf.Reset()
	}

	/*
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
	*/
	return nil
}

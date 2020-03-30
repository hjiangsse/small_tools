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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

type ShadowConfig struct {
	Svaddr   string `json:"server"`
	Lcaddr   string `json:"local_address"`
	Lcport   uint32 `json:"local_port"`
	Svport   uint32 `json:"server_port"`
	Password string `json:"password"`
	Timeout  uint32 `json:"timeout"`
	Method   string `json:"method"`
}

const (
	CONFPATHSTR string = "configpath"
	SVADDRSTR   string = "serveraddr"
	LCADDRSTR   string = "localaddr"
	LCPORTSTR   string = "localport"
	SVPORTSTR   string = "serverport"
	PASSWDSTR   string = "password"
	TIMEOUTSTR  string = "timeout"
	METHODSTR   string = "method"
)

// modCmd represents the mod command
var (
	configpath string
	serveraddr string
	localaddr  string
	localport  uint32
	serverport uint32
	password   string
	timeout    uint32
	method     string
)

var modCmd = &cobra.Command{
	Use:   "mod",
	Short: "modify shadowsock config infomation",
	Long:  `ask user input the new config information for shadowsocks`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := loadConfigInfo("/etc/shadowsocks.json")
		if err != nil {
			panic(err)
		}

		fmt.Println("Original Config Infomation:")
		fmt.Println("---------------------------------")
		fmt.Printf("%+v\n", cfg)
		fmt.Println("---------------------------------")

		changeConfigInfo(cmd, &cfg)
		fmt.Println("Changed Config Infomation:")
		fmt.Println("---------------------------------")
		fmt.Printf("%+v\n", cfg)
		fmt.Println("---------------------------------")

		//write new config info back to file
		writeConfBack(&cfg, "./temp_shadowsocks.json")

		//use sudo to replace the old config file
		execCmd := exec.Command("/bin/sh", "-c", "sudo mv ./temp_shadowsocks.json /etc/shadowsocks.json")
		execCmd.Run()

		psCmd := exec.Command("/bin/sh", "-c", "ps aux|grep sslocal|grep python")
		psOut, err := psCmd.Output()
		if err != nil {
			panic(err)
		}

		res := strings.Split(string(psOut), "\n")
		for _, line := range res {
			if strings.Contains(line, "shadowsocks.json") {
				sslocalPrcessId := strings.Fields(line)[1]
				killPara := fmt.Sprintf("sudo kill -9 %s", sslocalPrcessId)
				killCmd := exec.Command("/bin/sh", "-c", killPara)
				killCmd.Run()
			}
		}

		restartCmd := exec.Command("/bin/sh", "-c", "sslocal -c /etc/shadowsocks.json &")
		restartCmd.Run()
	},
}

func init() {
	rootCmd.AddCommand(modCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	modCmd.PersistentFlags().StringVarP(&configpath, CONFPATHSTR, "c", "/etc/shadowsocks.json", "config file path")
	modCmd.PersistentFlags().StringVarP(&serveraddr, SVADDRSTR, "s", "127.0.0.1", "the new server address you want to plant into the config")
	modCmd.PersistentFlags().StringVarP(&localaddr, LCADDRSTR, "l", "127.0.0.1", "the new local address you want to plant into the config")
	modCmd.PersistentFlags().Uint32VarP(&localport, LCPORTSTR, "p", 8080, "the new local port you want to plant into the config")
	modCmd.PersistentFlags().Uint32VarP(&serverport, SVPORTSTR, "v", 8080, "the new server port you want to plant into the config")
	modCmd.PersistentFlags().Uint32VarP(&timeout, TIMEOUTSTR, "t", 600, "the new timeout you want to plant into the config")
	modCmd.PersistentFlags().StringVarP(&method, METHODSTR, "m", "aes-256-cfb", "the new method you want to plant into the config")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// modCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

//load infomation from the config file, the file is json file(default path: /etc/shadowsocks.json)
func loadConfigInfo(path string) (ShadowConfig, error) {
	//read file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return ShadowConfig{}, err
	}

	var cfg ShadowConfig
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return ShadowConfig{}, err
	}

	return cfg, nil
}

func changeConfigInfo(cmd *cobra.Command, cfg *ShadowConfig) {
	if cmd.Flags().Changed(SVADDRSTR) {
		cfg.Svaddr = serveraddr
	}

	if cmd.Flags().Changed(LCADDRSTR) {
		cfg.Lcaddr = localaddr
	}

	if cmd.Flags().Changed(LCPORTSTR) {
		cfg.Lcport = localport
	}

	if cmd.Flags().Changed(SVPORTSTR) {
		cfg.Svport = serverport
	}

	if cmd.Flags().Changed(PASSWDSTR) {
		cfg.Password = password
	}

	if cmd.Flags().Changed(TIMEOUTSTR) {
		cfg.Timeout = timeout
	}

	if cmd.Flags().Changed(METHODSTR) {
		cfg.Method = method
	}
}

func writeConfBack(cfg *ShadowConfig, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	b, err := json.Marshal(*cfg)
	if err != nil {
		return err
	}

	file.Write(b)

	return nil
}

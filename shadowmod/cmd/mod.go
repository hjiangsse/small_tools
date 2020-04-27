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
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"strconv"
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

type DnsIpConfig struct {
	Dnshost string `json:"host"`
	IpAddr  string `json:"latest_address"`
}

type ProxyConfig struct {
	ProxyType      uint32 `json:"proxytype"`
	SocksRemoteDns bool   `json:"socksdns"`
	SocksAddr      string `json:"socksaddr"`
	SocksPort      uint32 `json:"socksport"`
}

const (
	PROXYPATHSTR string = "proxycnfpath"
	DNSPATHSTR   string = "dnscnfpath"
	CONFPATHSTR  string = "configpath"
	SVADDRSTR    string = "serveraddr"
	LCADDRSTR    string = "localaddr"
	LCPORTSTR    string = "localport"
	SVPORTSTR    string = "serverport"
	PASSWDSTR    string = "password"
	TIMEOUTSTR   string = "timeout"
	METHODSTR    string = "method"
)

// modCmd represents the mod command
var (
	proxycnfpath string
	dnscnfpath   string
	configpath   string
	serveraddr   string
	localaddr    string
	localport    uint32
	serverport   uint32
	password     string
	timeout      uint32
	method       string
)

var modCmd = &cobra.Command{
	Use:   "mod",
	Short: "modify shadowsock config infomation",
	Long:  `ask user input the new config information for shadowsocks`,
	Run: func(cmd *cobra.Command, args []string) {
		//load shadowsocks.json
		cfg, err := loadConfigInfo(configpath)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		//load shadowdns.json
		dnscfg, err := loadDnsConfInfo(dnscnfpath)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		//load shadowproxy.json
		proxycfg, err := loadProxyConfInfo(proxycnfpath)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		//fmt.Println(proxycfg)

		//after change the shadowsocks configs, change firefox profiles
		//and start a new instance of firfox for user
		err = doFirefoxProfileChange(&proxycfg)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		//compare latest dns ip address with the privious address
		dnsip, err := getLatestDnsIp(dnscfg.Dnshost)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		if dnscfg.IpAddr == dnsip {
			log.Printf("The previous dns ip is %v, the latest is %v, not change! Go on, boy!\n",
				dnscfg.IpAddr, dnsip)
			os.Exit(0)
		}

		//if dns ip changed, change the config files
		err = doCnfsChangeAndClientRestart(cmd, &cfg, &dnscfg, dnsip)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(modCmd)

	modCmd.PersistentFlags().StringVarP(&dnscnfpath, DNSPATHSTR, "d", "/etc/shadowdns.json", "dns config file path")
	modCmd.PersistentFlags().StringVarP(&proxycnfpath, PROXYPATHSTR, "x", "/etc/shadowproxy.json", "proxy config file path")
	modCmd.PersistentFlags().StringVarP(&configpath, CONFPATHSTR, "c", "/etc/shadowsocks.json", "config file path")
	modCmd.PersistentFlags().StringVarP(&serveraddr, SVADDRSTR, "s", "127.0.0.1", "the new server address you want to plant into the config")
	modCmd.PersistentFlags().StringVarP(&localaddr, LCADDRSTR, "l", "127.0.0.1", "the new local address you want to plant into the config")
	modCmd.PersistentFlags().Uint32VarP(&localport, LCPORTSTR, "p", 8080, "the new local port you want to plant into the config")
	modCmd.PersistentFlags().Uint32VarP(&serverport, SVPORTSTR, "v", 33128, "the new server port you want to plant into the config")
	modCmd.PersistentFlags().StringVarP(&password, PASSWDSTR, "a", "hZdHLzqdM3", "the new password you want to plant into the config")
	modCmd.PersistentFlags().Uint32VarP(&timeout, TIMEOUTSTR, "t", 600, "the new timeout you want to plant into the config")
	modCmd.PersistentFlags().StringVarP(&method, METHODSTR, "m", "aes-256-cfb", "the new method you want to plant into the config")
}

func doFirefoxProfileChange(cfg *ProxyConfig) error {
	//find firefox profile file
	curusr, err := user.Current()
	if err != nil {
		return err
	}

	//get the path of user preference file
	profilepath := curusr.HomeDir + "/.mozilla/firefox/*-release/prefs.js"
	profiledir := curusr.HomeDir + "/.mozilla/firefox/"

	//check and change proxy type
	newtype := "\"network.proxy.type\", " + strconv.Itoa(int(cfg.ProxyType)) + ");"
	newaddr := "\"network.proxy.socks\", " + "\"" + cfg.SocksAddr + "\"" + ");"
	newport := "\"network.proxy.socks_port\", " + strconv.Itoa(int(cfg.SocksPort)) + ");"
	newdns := "\"network.proxy.socks_remote_dns\", " + strconv.FormatBool(cfg.SocksRemoteDns) + ");"

	chgTypeStr := fmt.Sprintf("sed -i 's/\\(\"network.proxy.type\", .*\\)/%v/g' %v", newtype, profilepath)
	chgAddrStr := fmt.Sprintf("sed -i 's/\\(\"network.proxy.socks\", .*\\)/%v/g' %v", newaddr, profilepath)
	chgPortStr := fmt.Sprintf("sed -i 's/\\(\"network.proxy.socks_port\", .*\\)/%v/g' %v", newport, profilepath)
	chgDnsStr := fmt.Sprintf("sed -i 's/\\(\"network.proxy.socks_remote_dns\", .*\\)/%v/g' %v", newdns, profilepath)

	//do file change
	chgTypeCmd := exec.Command("/bin/sh", "-c", chgTypeStr)
	err = chgTypeCmd.Run()
	if err != nil {
		return err
	}

	chgAddrCmd := exec.Command("/bin/sh", "-c", chgAddrStr)
	err = chgAddrCmd.Run()
	if err != nil {
		return err
	}

	chgPortCmd := exec.Command("/bin/sh", "-c", chgPortStr)
	err = chgPortCmd.Run()
	if err != nil {
		return err
	}

	chgDnsCmd := exec.Command("/bin/sh", "-c", chgDnsStr)
	err = chgDnsCmd.Run()
	if err != nil {
		return err
	}

	//crate a new mozzila firefox instance
	err = os.Chdir(profiledir)
	if err != nil {
		return err
	}

	fds, err := ioutil.ReadDir(".")
	if err != nil {
		return nil
	}

	releasedir := profiledir
	for _, fd := range fds {
		if strings.Contains(fd.Name(), ("-release")) {
			releasedir = releasedir + fd.Name()
		}
	}

	profileabspath := releasedir + "/prefs.js"
	newInsCmdStr := fmt.Sprintf("firefox -P %v --new-instance &", profileabspath)
	newInsCmd := exec.Command("/bin/sh", "-c", newInsCmdStr)
	err = newInsCmd.Run()
	if err != nil {
		return err
	}

	fmt.Println("------------------------------------------------------")
	fmt.Println("-------外面的世界很精彩，外面的世界也很无奈-----------")
	fmt.Println("------------------------------------------------------")

	return nil
}

// change shadowsocks.json and shadowdns.json if latest dns ip changed
func doCnfsChangeAndClientRestart(cmd *cobra.Command, cnf *ShadowConfig, dnscnf *DnsIpConfig, latestip string) error {
	fmt.Println("Original Config Information:")
	fmt.Println("---------------------------------")
	fmt.Printf("%+v\n", *cnf)
	fmt.Println("---------------------------------")

	changeConfigInfo(cmd, cnf)
	cnf.Svaddr = latestip
	fmt.Println("Changed Config Information:")
	fmt.Println("---------------------------------")
	fmt.Printf("%+v\n", *cnf)
	fmt.Println("---------------------------------")

	//write new config info back to file
	err := writeCnfStructBackToFile(*cnf, "./temp_shadowsocks.json")
	if err != nil {
		return err
	}

	//write new dns config info back to file
	dnscnf.IpAddr = latestip
	err = writeCnfStructBackToFile(*dnscnf, "./temp_shadowdns.json")
	if err != nil {
		return err
	}

	//use sudo to replace the old config file
	mvCfgCmd := exec.Command("/bin/sh", "-c", "sudo mv ./temp_shadowsocks.json /etc/shadowsocks.json")
	err = mvCfgCmd.Run()
	if err != nil {
		return err
	}

	mvDnsCmd := exec.Command("/bin/sh", "-c", "sudo mv ./temp_shadowdns.json /etc/shadowdns.json")
	err = mvDnsCmd.Run()
	if err != nil {
		return err
	}

	psCmd := exec.Command("/bin/sh", "-c", "ps aux|grep sslocal|grep python")
	psOut, err := psCmd.Output()
	if err != nil {
		return err
	}

	res := strings.Split(string(psOut), "\n")
	for _, line := range res {
		if strings.Contains(line, "shadowsocks.json") {
			sslocalPrcessId := strings.Fields(line)[1]
			killPara := fmt.Sprintf("sudo kill -9 %s", sslocalPrcessId)
			killCmd := exec.Command("/bin/sh", "-c", killPara)

			err := killCmd.Run()
			if err != nil {
				return err
			}
		}
	}

	restartCmd := exec.Command("/bin/sh", "-c", "sslocal -c /etc/shadowsocks.json &")
	err = restartCmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// get the latest ip address from the dns host
func getLatestDnsIp(host string) (string, error) {
	nowaddr, err := net.LookupIP(host)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return nowaddr[0].String(), nil
}

// load shadowsocks proxy config info
func loadProxyConfInfo(path string) (ProxyConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return ProxyConfig{}, err
	}

	var cfg ProxyConfig
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return ProxyConfig{}, err
	}

	return cfg, nil
}

// load dns and ip config file, the file is a json file(default path: /etc/shadowdns.json)
func loadDnsConfInfo(path string) (DnsIpConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return DnsIpConfig{}, err
	}

	var cfg DnsIpConfig
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return DnsIpConfig{}, err
	}

	return cfg, nil
}

//load infomation from the config file, the file is json file(default path: /etc/shadowsocks.json)
func loadConfigInfo(path string) (ShadowConfig, error) {
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

func writeCnfStructBackToFile(v interface{}, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	file.Write(b)

	return nil
}

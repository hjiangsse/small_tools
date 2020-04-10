package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/mreiferson/go-options"
	"xchg.ai/sse/nsq/nsqd"
	"xchg.ai/sse/usensq/internal/app"
)

type program struct {
	once sync.Once
	nsqd *nsqd.NSQD
}

func (p *program) Start(configFile string) error {
	opts := nsqd.NewOptions()

	flagSet := app.NsqdFlagSet(opts)
	flagSet.Parse(os.Args[1:])

	rand.Seed(time.Now().UTC().UnixNano())

	if flagSet.Lookup("version").Value.(flag.Getter).Get().(bool) {
		fmt.Println(version.String("nsqd"))
		os.Exit(0)
	}

	var cfg config
	if configFile != "" {
		_, err := toml.DecodeFile(configFile, &cfg)
		if err != nil {
			logFatal("failed to load config file %s - %s", configFile, err)
		}
	}
	cfg.Validate()

	options.Resolve(opts, flagSet, cfg)
	nsqd, err := nsqd.New(opts)
	if err != nil {
		logFatal("failed to instantiate nsqd - %s", err)
	}
	p.nsqd = nsqd

	err = p.nsqd.LoadMetadata()
	if err != nil {
		logFatal("failed to load metadata - %s", err)
	}
	err = p.nsqd.PersistMetadata()
	if err != nil {
		logFatal("failed to persist metadata - %s", err)
	}

	go func() {
		err := p.nsqd.Main()
		if err != nil {
			p.Stop()
			os.Exit(1)
		}
	}()

	return nil
}

func (p *program) Stop() error {
	p.once.Do(func() {
		p.nsqd.Exit()
	})
	return nil
}

func main() {
	prg := &program{}
	if err := prg.Start(""); err != nil {
		fmt.Printf("NSQD Start Error: %v\n", err.Error())
	}

	sig = []os.Signal{syscall.SIGINT, syscall.SIGTERM}

	signalChan := make(chan os.Signal, 1)
	signalNotify(signalChan, sig...)
	<-signalChan

	if err := service.Stop(); err != nil {
		fmt.Printf("NSQD Stop Error: %v\n", err.Error())
	}
}

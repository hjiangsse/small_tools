package main

import (
	"os"
	"path/filepath"
	"sync"

	"example.com/myarch/internal/lg"
	"example.com/myarch/myarch"
	"github.com/judwhite/go-svc/svc"
)

type program struct {
	once   sync.Once
	myarch *myarch.MYARCH
}

func main() {
	prg := &program{}
	if err := svc.Run(prg); err != nil {
		logFatal("%s", err)
	}
}

func (p *program) Init(env svc.Environment) error {
	if env.IsWindowsService() {
		dir := filepath.Dir(os.Args[0])
		return os.Chdir(dir)
	}
	return nil

}

func (p *program) Start() error {
	//...Do some start prepare working
	go func() {
		err := p.myarch.Main()
		if err != nil {
			p.Stop()
			os.Exit(1)
		}
	}()
	return nil
}

func (p *program) Stop() error {
	p.once.Do(func() {
		p.myarch.Exit()
	})
	return nil
}

func logFatal(f string, args ...interface{}) {
	lg.LogFatal("[nsqd] %s", f, args...)
}

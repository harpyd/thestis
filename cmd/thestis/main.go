package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/harpyd/thestis/internal/runner"
)

const defaultConfigsPath = "configs/thestis"

func main() {
	flag.Parse()

	configsPath := flag.Arg(0)
	if configsPath == "" {
		configsPath = defaultConfigsPath
	}

	interrupted := make(chan os.Signal, 1)
	signal.Notify(interrupted, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	r := runner.New(configsPath)

	go r.Start()

	<-interrupted
	r.Stop()
}

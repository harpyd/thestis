package main

import (
	"flag"

	"github.com/harpyd/thestis/internal/runner"
)

const defaultConfigsPath = "configs/thestis"

func main() {
	flag.Parse()

	configsPath := flag.Arg(0)
	if configsPath == "" {
		configsPath = defaultConfigsPath
	}

	r := runner.New(configsPath)

	r.Start()
}

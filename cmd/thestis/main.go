package main

import (
	"flag"

	"github.com/harpyd/thestis/internal/runner"
)

const defaultConfigsPath = "configs"

func main() {
	flag.Parse()

	configsPath := flag.Arg(0)
	if configsPath == "" {
		configsPath = defaultConfigsPath
	}

	runner.Start(configsPath)
}

package main

import (
	"flag"

	"github.com/harpyd/thestis/internal/runner"
)

const defaultConfigsDir = "configs"

func main() {
	flag.Parse()

	configsDir := flag.Arg(0)
	if configsDir == "" {
		configsDir = defaultConfigsDir
	}

	runner.Start(configsDir)
}

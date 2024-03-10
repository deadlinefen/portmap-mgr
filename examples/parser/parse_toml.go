package main

import (
	"flag"
	"fmt"

	"github.com/deadlinefen/tinyPortMapper-manager-ipv6/pkg/config"
)

var tomlPath string

func init() {
	flag.StringVar(&tomlPath, "c", "config.toml", "toml config file path")
	flag.Parse()
}

func main() {
	parser := config.NewParserFactory(tomlPath).NewParser()

	config, err := parser.Parse()
	if err != nil {
		panic(fmt.Sprintf("%+v parser parse toml failed.", err))
	}

	fmt.Printf("config: %+v\n", config)
}

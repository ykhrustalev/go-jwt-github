package main

import (
	"flag"
	"github.com/ykhrustalev/exploregithub"
	"log"
)

var arguments struct {
	Config string
}

func init() {
	flag.StringVar(&arguments.Config, "config", "config.yml", "configuration file")
	flag.Parse()
}

func main() {
	if arguments.Config == "" {
		log.Fatalln("no config specified")
	}

	config, err := exploregithub.NewConfig(arguments.Config)
	if err != nil {
		log.Fatalf("can't read config, %v", err)
	}

	n := exploregithub.Server(config)
	n.Run(config.Server.Bind)
}

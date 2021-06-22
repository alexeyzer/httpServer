package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/alexeyzer/httpServer/internal/app/httpserver"
	"log"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "./configs/httpserver.toml", "path to config file")
}

func main() {

	flag.Parse()
	config := httpserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	server := httpserver.NewHttpServer(config)
	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}

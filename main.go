package main

import (
	"fs/src/cache"
	"fs/src/config"
	"fs/src/server"
	"github.com/spf13/pflag"
	"log"
)

var cfgPath string
var cfg *config.Config

func init() {
	pflag.StringVarP(&cfgPath, "config", "c", "./config.yml", "The path to the config file")
	pflag.Parse()

	var err error
	cfg, err = config.Load(cfgPath)
	if err != nil {
		log.Fatalln("error loading config.yml:", err)
	}

	err = cfg.Validate()
	if err != nil {
		log.Fatalln("invalid config:", err)
	}
}

func main() {
	cash := cache.New(cfg.Cache)

	server.Start(cfg, cash)
}

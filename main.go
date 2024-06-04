package main

import (
	"fs/src/cache"
	"fs/src/config"
	"fs/src/server"
	"log"
)

var cfg *config.Config

func init() {
	var err error
	cfg, err = config.Load("config.example.yml")
	if err != nil {
		log.Fatalln("error loading config.example.yml:", err)
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

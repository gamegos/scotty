package main

import (
	"flag"
	"log"
	"os"
	"runtime"

	"github.com/gamegos/scotty/config"
	"github.com/gamegos/scotty/server"
	"github.com/gamegos/scotty/storage"
	//_ "github.com/gamegos/scotty/storage/drivers/memory"
	_ "github.com/gamegos/scotty/storage/drivers/redis"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	confPath := flag.String("config", "", "Config file")
	flag.Parse()

	var conf *config.Config
	if *confPath != "" {
		confFile, err := os.Open(*confPath)
		defer confFile.Close()
		if err != nil {
			log.Fatalf("could not load config file: %s, err: %s", confPath, err)
		}

		c, err := config.Parse(confFile)
		if err != nil {
			log.Fatalf("could not parse config %s", err)
		}
		conf = c
	} else {
		log.Println("using default config")
		conf = config.DefaultConfig()
	}

	stg := storage.Init(conf.Storage.Driver, conf.Storage.Options)

	log.Printf("starting scotty server on %s", conf.Server.Addr)
	s := server.Init(stg)
	log.Fatal(s.Run(conf.Server.Addr))
}

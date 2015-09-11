package main

import (
	"flag"
	"log"
	"runtime"

	"github.com/gamegos/scotty/server"
	"github.com/gamegos/scotty/storage"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	confFile := flag.String("config", "default.conf", "Config file")
	flag.Parse()

	conf := storage.InitConfig(*confFile)
	stg := storage.Init(&conf.Redis)

	log.Printf("starting scotty server on %s", conf.Addr)
	s := server.Init(stg)
	log.Fatal(s.Run(conf.Addr))
}

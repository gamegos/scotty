package main

import (
	"flag"
	"log"
	"net/http"
	"runtime"

	"github.com/gamegos/scotty/service"
	"github.com/gamegos/scotty/storage"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	confFile := flag.String("config", "default.conf", "Config file")
	flag.Parse()

	conf := storage.InitConfig(*confFile)
	stg := storage.Init(&conf.Redis)

	router := service.NewRouter(stg)
	log.Printf("scotty service listening on %s", conf.Addr)
	log.Fatal(http.ListenAndServe(conf.Addr, router))
}

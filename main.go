package main

import (
	"flag"
	"log"
	"runtime"
	"net/http"

	"gitlab.fixb.com/mir/push/storage"
	"gitlab.fixb.com/mir/push/service"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	confFile := flag.String("config", "", "Config file")
	flag.Parse()

	conf := storage.InitConfig(*confFile)
	stg := storage.Init(&conf.Redis)

	router := service.NewRouter(stg)
	log.Fatal(http.ListenAndServe(conf.Addr, router))
}
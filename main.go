package main

import (
	"flag"
	"log"
	"net/http"
	"runtime"

	"gitlab.fixb.com/mir/push/service"
	"gitlab.fixb.com/mir/push/storage"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	confFile := flag.String("config", "default.conf", "Config file")
	flag.Parse()

	conf := storage.InitConfig(*confFile)
	stg := storage.Init(&conf.Redis)

	router := service.NewRouter(stg)
	log.Printf("Push service listening on %s", conf.Addr)
	log.Fatal(http.ListenAndServe(conf.Addr, router))
}

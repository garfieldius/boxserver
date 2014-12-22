package main

import (
	"github.com/trenker/boxserver/data"
	"github.com/trenker/boxserver/log"
	"github.com/trenker/boxserver/conf"
	"github.com/trenker/boxserver/server"
	"flag"
	"net/http"
)

func main() {

	log.Debug("Load config")

	var configSrc string

	flag.StringVar(&configSrc, "c", "default", "")
	flag.Parse()

	conf.Load(configSrc)

	log.Debug("Loading data")
	data.Initialize(conf.Get().Data)

	log.Debug("Register request handler")
	http.Handle("/", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		server.NewRequest(req).Process(res)
	}))


	port := conf.Get().Port
	log.Debug("Listen on port %s", port)

	err := http.ListenAndServe(":" + port, nil)

	if err != nil {
		log.Critical("Cannot listen to requests, %s", err)
	}
}

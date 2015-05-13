package main

import (
	"flag"
	"github.com/trenker/boxserver/conf"
	"github.com/trenker/boxserver/data"
	"github.com/trenker/boxserver/log"
	"github.com/trenker/boxserver/server"
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

	if conf.Get().Server.Enable {
		log.Debug("Register fileserver of route %s for directory %s", conf.Get().Server.Prefix, conf.Get().Server.BaseDir)
		http.Handle(conf.Get().Server.Prefix, http.StripPrefix(conf.Get().Server.Prefix, http.FileServer(http.Dir(conf.Get().Server.BaseDir))))
	}

	log.Debug("Register request handler")
	http.Handle("/", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		server.NewRequest(req).Process(res)
	}))

	port := conf.Get().Port
	log.Debug("Listen on port %s", port)

	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Critical("Cannot listen to requests, %s", err)
	}
}

package main

import (
	"flag"
	"log"

	appmod "quickChain/app"

	"github.com/cometbft/cometbft/abci/server"
)

func main() {
	demo := flag.Bool("demo", false, "run demo client instead of server")
	flag.Parse()

	if *demo {
		if err := runDemo(); err != nil {
			log.Fatal(err)
		}
		return
	}

	application := appmod.NewDataStoreApp()
	srv, err := server.NewServer("tcp://127.0.0.1:26658", "socket", application)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Server started")
	}
	srv.Start()
	defer srv.Stop()
	select {} // block forever
}

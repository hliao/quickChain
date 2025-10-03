package main

import (
	"flag"
	"log"

	appmod "quickChain/app"

	"github.com/cometbft/cometbft/abci/server"
)

func main() {
	demo := flag.Bool("demo", false, "run demo client instead of server")
	transport := flag.String("transport", "socket", "ABCI transport: socket|grpc")
	addr := flag.String("addr", "tcp://127.0.0.1:26658", "ABCI listen address")
	flag.Parse()

	if *demo {
		if err := runDemo(); err != nil {
			log.Fatal(err)
		}
		return
	}

	application := appmod.NewDataStoreApp()
	srv, err := server.NewServer(*addr, *transport, application)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Server started on %s (%s)", *addr, *transport)
	}
	srv.Start()
	defer srv.Stop()
	select {} // block forever
}

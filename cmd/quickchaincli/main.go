package main

import (
	"flag"
	"log"

	appmod "quickChain/app"

	"github.com/cometbft/cometbft/abci/server"
)

func main() {
	transport := flag.String("transport", "socket", "ABCI transport: socket|grpc")
	addr := flag.String("addr", "tcp://127.0.0.1:26658", "ABCI listen address")
	flag.Parse()

	application := appmod.NewDataStoreApp()
	srv, err := server.NewServer(*addr, *transport, application)
	if err != nil {
		log.Fatal(err)
	}
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
	log.Printf("Server started on %s (%s)", *addr, *transport)
	defer srv.Stop()
	<-srv.Quit()
}

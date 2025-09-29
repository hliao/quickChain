package main

import (
	"log"

	appmod "quickChain/app"

	"github.com/cometbft/cometbft/abci/server"
)

func main() {
	application := &appmod.CounterApp{}
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

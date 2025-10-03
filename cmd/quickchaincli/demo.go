package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	abcicli "github.com/cometbft/cometbft/abci/client"
	abci "github.com/cometbft/cometbft/abci/types"
)

func runDemo() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := abcicli.NewClient("tcp://127.0.0.1:26658", "socket", true)
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}
	if err := client.Start(); err != nil {
		return fmt.Errorf("start client: %w", err)
	}
	defer func() {
		if stopErr := client.Stop(); stopErr != nil {
			log.Printf("client stop error: %v", stopErr)
		}
	}()

	// generate 1024-byte random chunk
	tx := make([]byte, 1024)
	if _, err := rand.Read(tx); err != nil {
		return fmt.Errorf("rand.Read: %w", err)
	}

	// optional: CheckTx
	if _, err := client.CheckTx(ctx, &abci.CheckTxRequest{Tx: tx}); err != nil {
		return fmt.Errorf("CheckTx: %w", err)
	}

	// FinalizeBlock with one tx
	fb, err := client.FinalizeBlock(ctx, &abci.FinalizeBlockRequest{Txs: [][]byte{tx}})
	if err != nil {
		return fmt.Errorf("FinalizeBlock: %w", err)
	}
	if len(fb.TxResults) != 1 || fb.TxResults[0].Code != abci.CodeTypeOK {
		return fmt.Errorf("FinalizeBlock not OK: %+v", fb.TxResults)
	}
	key := string(fb.TxResults[0].Data)
	if _, err := hex.DecodeString(key); err != nil {
		return fmt.Errorf("invalid key hex: %q", key)
	}
	log.Printf("Stored chunk with key: %s", key)

	// Query it back
	q, err := client.Query(ctx, &abci.QueryRequest{Path: "/get", Data: []byte(key)})
	if err != nil {
		return fmt.Errorf("Query: %w", err)
	}
	if q.Code != abci.CodeTypeOK {
		return fmt.Errorf("query not OK: code=%d log=%s", q.Code, q.Log)
	}
	if len(q.Value) != len(tx) {
		return fmt.Errorf("size mismatch: got %d want %d", len(q.Value), len(tx))
	}
	log.Printf("Retrieved %d bytes successfully", len(q.Value))
	return nil
}

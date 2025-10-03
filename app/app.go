package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	abci "github.com/cometbft/cometbft/abci/types"
)

const chunkSize = 1024

type DataStoreApp struct {
	abci.BaseApplication
	store map[string][]byte
}

func NewDataStoreApp() *DataStoreApp {
	return &DataStoreApp{store: make(map[string][]byte)}
}

// FinalizeBlock: persist each 1024-byte tx and return its hash key in Data
func (app *DataStoreApp) FinalizeBlock(_ context.Context, req *abci.FinalizeBlockRequest) (*abci.FinalizeBlockResponse, error) {
	if app.store == nil {
		app.store = make(map[string][]byte)
	}
	results := make([]*abci.ExecTxResult, len(req.Txs))
	for i, tx := range req.Txs {
		if len(tx) != chunkSize {
			msg := fmt.Sprintf("invalid chunk size: got %d want %d", len(tx), chunkSize)
			log.Printf("FinalizeBlock reject: %s", msg)
			results[i] = &abci.ExecTxResult{Code: 1, Log: msg}
			continue
		}
		sum := sha256.Sum256(tx)
		key := hex.EncodeToString(sum[:])
		data := make([]byte, len(tx))
		copy(data, tx)
		app.store[key] = data
		log.Printf("FinalizeBlock stored: key=%s size=%d", key, len(tx))
		results[i] = &abci.ExecTxResult{Code: abci.CodeTypeOK, Data: []byte(key)}
	}
	return &abci.FinalizeBlockResponse{TxResults: results}, nil
}

// CheckTx: validate size before entering mempool
func (app *DataStoreApp) CheckTx(_ context.Context, req *abci.CheckTxRequest) (*abci.CheckTxResponse, error) {
	if len(req.Tx) != chunkSize {
		msg := fmt.Sprintf("invalid chunk size: got %d want %d", len(req.Tx), chunkSize)
		log.Printf("CheckTx reject: %s", msg)
		return &abci.CheckTxResponse{Code: 1, Log: msg}, nil
	}
	log.Printf("CheckTx ok: size=%d", len(req.Tx))
	return &abci.CheckTxResponse{Code: abci.CodeTypeOK}, nil
}

// Query: path "get" with Data=hex(hash) returns the stored 1024-byte value
func (app *DataStoreApp) Query(_ context.Context, req *abci.QueryRequest) (*abci.QueryResponse, error) {
	if app.store == nil {
		app.store = make(map[string][]byte)
	}
	switch strings.Trim(req.Path, "/") {
	case "get":
		key := strings.TrimSpace(string(req.Data))
		if key == "" {
			log.Printf("Query get: empty key")
			return &abci.QueryResponse{Code: 1, Log: "empty key"}, nil
		}
		if val, ok := app.store[key]; ok {
			log.Printf("Query get: hit key=%s size=%d", key, len(val))
			return &abci.QueryResponse{Code: abci.CodeTypeOK, Value: val}, nil
		}
		log.Printf("Query get: miss key=%s", key)
		return &abci.QueryResponse{Code: 1, Log: "not found"}, nil
	default:
		log.Printf("Query unknown path: %q", req.Path)
		return &abci.QueryResponse{Code: 1, Log: "unknown path"}, nil
	}
}

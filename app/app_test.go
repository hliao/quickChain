package app

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
)

func randomBytes(t *testing.T, n int) []byte {
	t.Helper()
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		t.Fatalf("rand.Read: %v", err)
	}
	return b
}

func TestCheckTx_ValidAndInvalidSizes(t *testing.T) {
	app := NewDataStoreApp()
	// valid
	resp, err := app.CheckTx(context.Background(), &abci.CheckTxRequest{Tx: make([]byte, chunkSize)})
	if err != nil {
		t.Fatalf("CheckTx error: %v", err)
	}
	if resp.Code != abci.CodeTypeOK {
		t.Fatalf("expected OK, got %d (%s)", resp.Code, resp.Log)
	}
	// invalid
	resp, err = app.CheckTx(context.Background(), &abci.CheckTxRequest{Tx: make([]byte, chunkSize-1)})
	if err != nil {
		t.Fatalf("CheckTx error: %v", err)
	}
	if resp.Code == abci.CodeTypeOK {
		t.Fatalf("expected non-OK for invalid size")
	}
}

func TestFinalizeBlock_StoresChunkAndReturnsKey(t *testing.T) {
	app := NewDataStoreApp()
	tx := randomBytes(t, chunkSize)
	fbResp, err := app.FinalizeBlock(context.Background(), &abci.FinalizeBlockRequest{Txs: [][]byte{tx}})
	if err != nil {
		t.Fatalf("FinalizeBlock error: %v", err)
	}
	if len(fbResp.TxResults) != 1 {
		t.Fatalf("expected 1 result, got %d", len(fbResp.TxResults))
	}
	res := fbResp.TxResults[0]
	if res.Code != abci.CodeTypeOK {
		t.Fatalf("expected OK, got %d (%s)", res.Code, res.Log)
	}
	key := string(res.Data)
	if _, err := hex.DecodeString(key); err != nil {
		t.Fatalf("invalid hex key: %q", key)
	}
	qResp, err := app.Query(context.Background(), &abci.QueryRequest{Path: "/get", Data: []byte(key)})
	if err != nil {
		t.Fatalf("Query error: %v", err)
	}
	if qResp.Code != abci.CodeTypeOK {
		t.Fatalf("query not OK: %d (%s)", qResp.Code, qResp.Log)
	}
	if len(qResp.Value) != chunkSize {
		t.Fatalf("unexpected value size: %d", len(qResp.Value))
	}
}

func TestQuery_NotFoundAndUnknownPath(t *testing.T) {
	app := NewDataStoreApp()
	// not found
	qResp, err := app.Query(context.Background(), &abci.QueryRequest{Path: "/get", Data: []byte("deadbeef")})
	if err != nil {
		t.Fatalf("Query error: %v", err)
	}
	if qResp.Code == abci.CodeTypeOK {
		t.Fatalf("expected not found")
	}
	// unknown path
	qResp, err = app.Query(context.Background(), &abci.QueryRequest{Path: "/unknown", Data: nil})
	if err != nil {
		t.Fatalf("Query error: %v", err)
	}
	if qResp.Code == abci.CodeTypeOK {
		t.Fatalf("expected error for unknown path")
	}
}

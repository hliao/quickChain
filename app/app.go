package app

import (
	"context"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
)

type CounterApp struct {
	abci.BaseApplication
	count int64
}

// FinalizeBlock: called with the decided block and its txs
func (app *CounterApp) FinalizeBlock(_ context.Context, req *abci.FinalizeBlockRequest) (*abci.FinalizeBlockResponse, error) {
	app.count += int64(len(req.Txs))
	results := make([]*abci.ExecTxResult, len(req.Txs))
	for i := range results {
		results[i] = &abci.ExecTxResult{Code: abci.CodeTypeOK}
	}
	return &abci.FinalizeBlockResponse{TxResults: results}, nil
}

// Query: allow querying current count
func (app *CounterApp) Query(_ context.Context, _ *abci.QueryRequest) (*abci.QueryResponse, error) {
	return &abci.QueryResponse{Code: abci.CodeTypeOK, Value: []byte(fmt.Sprintf("%d", app.count))}, nil
}

package types

import (
	"errors"
	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/state"
	abci "github.com/tendermint/tendermint/abci/types"
)

type BlockResult struct {
	Block *types.Block
	Results *state.ABCIResponses 
}

type TxResult struct {
	Tx types.Tx
	Result abci.ResponseDeliverTx
}

func (result *BlockResult) GetTx(index int64) (TxResult, error){
	if result.Block.NumTxs <= index || int64(len(result.Block.Txs)) <= index || int64(len(result.Results.DeliverTx)) <= index {
		return TxResult{}, errors.New("Out of range")
	}

	return TxResult {
		Tx: result.Block.Txs[index],
		Result: *result.Results.DeliverTx[index],
	}, nil
}

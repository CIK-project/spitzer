package types

import (
	"encoding/hex"
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/codec"
)

type Tx struct {
	Height    int64
	Index     int64
	Hash      string
	Code      int32
	Signers   []string
	Memo      string
	Tags      string
	Fee       string
	GasWanted int64
	GasUsed   int64
	Msgs      string
}

func NewTxFrom(txResult *TxResult, height int64, index int64, txDecoder sdk.TxDecoder, cdc *codec.Codec) (*Tx, error) {
	var stdTx auth.StdTx
	appTx, sdkErr := txDecoder(txResult.Tx)
	// if err occured, just use default stdTx
	if sdkErr != nil {
		stdTx = auth.StdTx{}
	} else {
		_stdTx, ok := appTx.(auth.StdTx)
		if ok == false {
			// return nil, errors.New("Tx is not std tx")
			stdTx = auth.StdTx{}
		} else {
			stdTx = _stdTx
		}
	}

	signers := make([]string, 0)
	for _, signature := range stdTx.GetSignatures() {
		signers = append(signers, sdk.AccAddress(signature.Address()).String())
	}

	memo := stdTx.GetMemo()

	tags := make([]map[string]string, 0)
	for _, tag := range txResult.Result.Tags {
		t := make(map[string]string)
		t[string(tag.Key)] = string(tag.Value)
		tags = append(tags, t)
	}
	tagsJson, err := json.Marshal(tags)
	if err != nil {
		return nil, err
	}

	fee := stdTx.Fee
	feeJson, err := json.Marshal(fee.Amount)
	if err != nil {
		return nil, err
	}

	msgsJson, err := cdc.MarshalJSON(stdTx.GetMsgs())
	if err != nil {
		return nil, err
	}

	return &Tx{
		Height:    height,
		Index:     index,
		Hash:      hex.EncodeToString(txResult.Tx.Hash()),
		Code:      int32(txResult.Result.Code),
		Signers:   signers,
		Memo:      memo,
		Tags:      string(tagsJson),
		Fee:       string(feeJson),
		GasWanted: txResult.Result.GasWanted,
		GasUsed:   txResult.Result.GasUsed,
		Msgs:      string(msgsJson),
	}, nil
}

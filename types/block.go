package types

import (
	"encoding/json"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Header struct {
	Height            int64
	Hash              string
	PrevHash          string
	Time              time.Time
	NumTxs            int64
	TotalTxs          int64
	LastCommitHash    string
	DataHash          string
	ValidatorHash     string
	NextValidatorHash string
	ConsensusHash     string
	AppHash           string
	LastResultHash    string
	EvidenceHash      string
	Proposer          string
	Tags              string
}

func NewHeaderFrom(result *BlockResult) *Header {
	tags := make([]map[string]string, 0)
	for _, tag := range result.Results.BeginBlock.Tags {
		if string(tag.Key) == "height" {
			continue
		}
		t := make(map[string]string)
		t[string(tag.Key)] = string(tag.Value)
		tags = append(tags, t)
	}
	for _, tag := range result.Results.EndBlock.Tags {
		if string(tag.Key) == "height" {
			continue
		}
		t := make(map[string]string)
		t[string(tag.Key)] = string(tag.Value)
		tags = append(tags, t)
	}
	tagsJson, err := json.Marshal(tags)
	if err != nil {
		panic("Fail to encode tags by json")
	}

	block := result.Block
	return &Header{
		Height:            block.Height,
		Hash:              block.Hash().String(),
		PrevHash:          block.LastBlockID.Hash.String(),
		Time:              block.Time,
		NumTxs:            block.NumTxs,
		TotalTxs:          block.TotalTxs,
		LastCommitHash:    block.LastCommitHash.String(),
		DataHash:          block.DataHash.String(),
		ValidatorHash:     block.ValidatorsHash.String(),
		NextValidatorHash: block.NextValidatorsHash.String(),
		ConsensusHash:     block.ConsensusHash.String(),
		AppHash:           block.AppHash.String(),
		LastResultHash:    block.LastResultsHash.String(),
		EvidenceHash:      block.EvidenceHash.String(),
		Proposer:          sdk.ConsAddress(block.ProposerAddress).String(),
		Tags:              string(tagsJson),
	}
}

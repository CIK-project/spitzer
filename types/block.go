package types

type Header struct {
	Height            int64
	Hash              string
	PrevHash          string
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
}

func NewHeaderFrom(result *BlockResult) *Header {
	block := result.Block
	return &Header{
		Height:            block.Height,
		Hash:              block.Hash().String(),
		PrevHash:          block.LastBlockID.Hash.String(),
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
		Proposer:          block.ProposerAddress.String(),
	}
}

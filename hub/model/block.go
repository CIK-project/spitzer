package model

import (
	"database/sql"
	"github.com/CIK-project/spitzer/types"
)

func GetHeader(db *sql.DB, height int64) (*types.Header, error) {
	header := types.Header{}
	row := db.QueryRow(`
	SELECT
		height,
		hash,
		prevHash,
		numTxs,
		totalTxs,
		lastCommitHash,
		dataHash,
		validatorHash,
		nextValidatorHash,
		consensusHash,
		appHash,
		lastResultHash,
		evidenceHash,
		proposer
	FROM block
	WHERE height = $1
	`, height)
	err := row.Scan(
		&header.Height,
		&header.Hash,
		&header.PrevHash,
		&header.NumTxs,
		&header.TotalTxs,
		&header.LastCommitHash,
		&header.DataHash,
		&header.ValidatorHash,
		&header.NextValidatorHash,
		&header.ConsensusHash,
		&header.AppHash,
		&header.LastResultHash,
		&header.EvidenceHash,
		&header.Proposer,
	)
	if err != nil {
		return nil, err
	}

	return &header, err
}
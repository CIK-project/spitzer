package model

import (
	"errors"
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
		time,
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
		proposer,
		tags
	FROM block
	WHERE height = $1
	`, height)
	err := row.Scan(
		&header.Height,
		&header.Hash,
		&header.PrevHash,
		&header.Time,
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
		&header.Tags,
	)
	if err != nil {
		return nil, err
	}

	return &header, err
}

func GetHeaders(db *sql.DB, page int, limit int) ([]*types.Header, error) {
	if page <= 0 {
		return nil, errors.New("Invalid page")
	}
	if limit <= 0 {
		return nil, errors.New("Invalid limit")
	}

	rows, err := db.Query(`
	SELECT
		height,
		hash,
		prevHash,
		time,
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
		proposer,
		tags
	FROM block
	ORDER BY height DESC
	LIMIT $1 OFFSET $2
	`, limit, (page - 1) * limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	headers := make([]*types.Header, 0, 10)
	
	for rows.Next() {
		header := types.Header{}
		err := rows.Scan(
			&header.Height,
			&header.Hash,
			&header.PrevHash,
			&header.Time,
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
			&header.Tags,
		)
		if err != nil {
			return nil, err
		}
		headers = append(headers, &header)
	}
	return headers, nil
}

func SetHeader(tx *sql.Tx, header *types.Header) (sql.Result, error) {
	return tx.Exec(`
		INSERT INTO block(
			height,
			hash,
			prevHash,
			time,
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
			proposer,
			tags
		)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		`, header.Height,
		header.Hash,
		header.PrevHash,
		header.Time,
		header.NumTxs,
		header.TotalTxs,
		header.LastCommitHash,
		header.DataHash,
		header.ValidatorHash,
		header.NextValidatorHash,
		header.ConsensusHash,
		header.AppHash,
		header.LastResultHash,
		header.EvidenceHash,
		header.Proposer,
		header.Tags,
	)
}

func RecentHeight(db *sql.DB) (int64, error) {
	row := db.QueryRow(`
	SELECT
		height
	FROM
		block
	ORDER BY
		height DESC
	LIMIT 1		
	`)
	var recentHeight int64
	err := row.Scan(&recentHeight)
	if err != nil {
		return 0, err
	}
	return recentHeight, nil
}

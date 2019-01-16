package model

import (
	"errors"
	"database/sql"

	"github.com/lib/pq"

	"github.com/CIK-project/spitzer/types"
)

func SetTx(dbTx *sql.Tx, tx *types.Tx) (sql.Result, error) {
	return dbTx.Exec(`
	INSERT INTO transaction(
		height,
		index,
		hash,
		code,
		signers,
		memo,
		tags,
		fee,
		gasWanted,
		gasUsed,
		msgs
	)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, tx.Height,
	tx.Index,
	tx.Hash,
	tx.Code,
	pq.Array(tx.Signers),
	tx.Memo,
	tx.Tags,
	tx.Fee,
	tx.GasWanted,
	tx.GasUsed,
	tx.Msgs)
}

func GetTxs(db *sql.DB, page int, limit int, height int64, signers string, tags string) ([]*types.Tx, error) {
	if page <= 0 {
		return nil, errors.New("Invalid page")
	}
	if limit <= 0 {
		return nil, errors.New("Invalid limit")
	}

	q := `
	SELECT
		height,
		index,
		hash,
		code,
		signers,
		memo,
		tags,
		fee,
		gasWanted,
		gasUsed,
		msgs
	FROM transaction
	WHERE
		tags @> $3
	`
	if len(signers) > 2 {
		q += "AND signers <@ $4"
	}
	if height > 0 {
		if len(signers) > 2 {
			q += "AND height = $5"
		}	else {
			q += "AND height = $4"
		}
	}
	q += `
	ORDER BY
		height DESC,
		index DESC
	LIMIT $1 OFFSET $2
	`
	
	var rows *sql.Rows
	var err error
	if len(signers) > 2 {
		if height > 0 {
			rows, err = db.Query(q, limit, (page - 1) * limit, tags, signers, height)
		} else {
			rows, err = db.Query(q, limit, (page - 1) * limit, tags, signers)
		}
	}	else {
		if height > 0 {
			rows, err = db.Query(q, limit, (page - 1) * limit, tags, height)
		} else {
			rows, err = db.Query(q, limit, (page - 1) * limit, tags)
		}
	}
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	txs := make([]*types.Tx, 0, 10)
	
	for rows.Next() {
		tx := types.Tx{}
		err := rows.Scan(
			&tx.Height,
			&tx.Index,
			&tx.Hash,
			&tx.Code,
			pq.Array(&tx.Signers),
			&tx.Memo,
			&tx.Tags,
			&tx.Fee,
			&tx.GasWanted,
			&tx.GasUsed,
			&tx.Msgs,
		)
		if err != nil {
			return nil, err
		}
		txs = append(txs, &tx)
	}
	return txs, nil
}
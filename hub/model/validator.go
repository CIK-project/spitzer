package model

import (
	"errors"
	"database/sql"
	"github.com/lib/pq"
	"github.com/CIK-project/spitzer/types"
)

func GetValidator(db *sql.DB, address string) (*types.Validator, error) {
	validator := types.Validator{}
	uptime := make([]int64, 0)
	row := db.QueryRow(`
	SELECT
		address,
		consAddress,
		pubkey,
		power,
		moniker,
    identity,
    website,
		details,
		uptime
	FROM validator
	WHERE address = $1
	ORDER BY id DESC
	LIMIT 1
	`, address)
	err := row.Scan(
		&validator.Address,
		&validator.ConsAddress,
		&validator.PubKey,
		&validator.Power,
		&validator.Moniker,
		&validator.Identity,
		&validator.Website,
		&validator.Details,
		pq.Array(&uptime),
	)
	if err != nil {
		return nil, err
	}
	var sum int64 = 0
	for _, u := range uptime {
		sum += u
	}
	validator.Uptime = int(sum)

	return &validator, err
}

func GetValidators(db *sql.DB, page int, limit int) ([]*types.Validator, error) {
	if page <= 0 {
		return nil, errors.New("Invalid page")
	}
	if limit <= 0 {
		return nil, errors.New("Invalid limit")
	}

	rows, err := db.Query(`
	SELECT
		address,
		consAddress,
		pubkey,
		power,
		moniker,
		identity,
		website,
		details,
		uptime
	FROM validator
	WHERE
		power > 0
	ORDER BY
		power DESC,
		id DESC
	LIMIT $1 OFFSET $2
	`, limit, (page - 1) * limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	validators := make([]*types.Validator, 0, 10)
	
	for rows.Next() {
		validator := types.Validator{}
		uptime := make([]int64, 0)
		err := rows.Scan(
			&validator.Address,
			&validator.ConsAddress,
			&validator.PubKey,
			&validator.Power,
			&validator.Moniker,
			&validator.Identity,
			&validator.Website,
			&validator.Details,
			pq.Array(&uptime),
		)
		if err != nil {
			return nil, err
		}
		var sum int64 = 0
		for _, u := range uptime {
			sum += u
		}
		validator.Uptime = int(sum)
		validators = append(validators, &validator)
	}
	return validators, nil
}

func GetValidatorByCons(db *sql.DB, consAddress string) (*types.Validator, error) {
	validator := types.Validator{}
	uptime := make([]int64, 0)
	row := db.QueryRow(`
	SELECT
		address,
		consAddress,
		pubkey,
		power,
		moniker,
    identity,
    website,
		details,
		uptime
	FROM validator
	WHERE consAddress = $1
	ORDER BY id DESC
	LIMIT 1
	`, consAddress)
	err := row.Scan(
		&validator.Address,
		&validator.ConsAddress,
		&validator.PubKey,
		&validator.Power,
		&validator.Moniker,
		&validator.Identity,
		&validator.Website,
		&validator.Details,
		pq.Array(&uptime),
	)
	if err != nil {
		return nil, err
	}
	var sum int64 = 0
	for _, u := range uptime {
		sum += u
	}
	validator.Uptime = int(sum)

	return &validator, err
}

func SetValidator(tx *sql.Tx, validator *types.Validator) (sql.Result, error) {
	return tx.Exec(`
		INSERT INTO validator(
			address,
			consAddress,
			pubkey,
			power
		)
		VALUES
			($1, $2, $3, $4)
		ON CONFLICT ON CONSTRAINT validator_pkey
		DO
			UPDATE
				SET power = $4	
		`, validator.Address,
		validator.ConsAddress,
		validator.PubKey,
		validator.Power,
	)
}

func SetValidatorPower(tx *sql.Tx, validator *types.Validator) (sql.Result, error) {
	return tx.Exec(`
	UPDATE validator
		SET power = $1
	WHERE
		id = (SELECT id FROM validator WHERE pubkey = $2 ORDER BY ID DESC LIMIT 1)
	`, validator.Power,
	validator.PubKey)
}

func SetValidatorDescription(tx *sql.Tx, address, moniker, identity, website, details string) (sql.Result, error) {
	return tx.Exec(`
		UPDATE validator
			SET moniker = $1,
					identity = $2,
					website = $3,
					details = $4
		WHERE
			id = (SELECT id FROM validator WHERE address = $5 ORDER BY ID DESC LIMIT 1)
	`, moniker,
	identity,
	website,
	details,
	address)
}

func DecrementValsUptime(tx *sql.Tx) (sql.Result, error) {
	return tx.Exec(`
		CALL decrementAllUptime();	
	`)
}

func IncrementValUptime(tx *sql.Tx, validator *types.Validator) (sql.Result, error) {
	return tx.Exec(`
		CALL increaseUptime($1);
	`, validator.ConsAddress)
}

func LimitValsUptimeLen(tx *sql.Tx) (sql.Result, error) {
	return tx.Exec(`
		CALL limitUptimeLen();	
	`)
}

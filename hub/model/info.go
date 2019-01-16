package model

import (
	"database/sql"
)

func GetNumOfValidators(db *sql.DB) (int, error) {
	count := 0
	row := db.QueryRow(`
	SELECT
		count(id)
	FROM validator
	WHERE power > 0
	`)
	err := row.Scan(
		&count,
	)
	if err != nil {
		return 0, err
	}

	return count, err
}

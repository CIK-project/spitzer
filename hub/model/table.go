package model

import (
	"database/sql"
)

func CreateTable(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS block
	(
		height BIGSERIAL PRIMARY KEY,
		hash CHAR (64) NOT NULL UNIQUE,
		prevHash CHAR (64) UNIQUE,
		time TIMESTAMP,
		numTxs BIGINT,
		totalTxs BIGINT,
		lastCommitHash CHAR (64),
		dataHash CHAR (64),
		validatorHash CHAR (64),
		nextValidatorHash CHAR (64),
		consensusHash CHAR (64),
		appHash CHAR (64),
		lastResultHash CHAR (64),
		evidenceHash CHAR (64),
		proposer TEXT,
		tags JSONB
	)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
	CREATE INDEX IF NOT EXISTS idx_block_time
	ON block (time)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
	CREATE INDEX IF NOT EXISTS idx_block_proposer
	ON block
	USING HASH (proposer)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
	CREATE INDEX IF NOT EXISTS idx_block_gin_tags
	ON block
	USING GIN (tags)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS validator
	(
		id SERIAL UNIQUE,
		address TEXT,
		consAddress TEXT UNIQUE,
		pubkey TEXT,
		power BIGINT,
		moniker TEXT DEFAULT '',
		identity TEXT DEFAULT '',
		website TEXT DEFAULT '',
		details TEXT DEFAULT '',
		uptime INTEGER[] DEFAULT '{}',
		PRIMARY KEY (address, pubkey)
	)`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`
	CREATE INDEX IF NOT EXISTS idx_validator_power
	ON validator (power)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
	CREATE INDEX IF NOT EXISTS idx_validator_moniker
	ON validator
	USING HASH (moniker)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
	CREATE OR REPLACE PROCEDURE decrementAllUptime()
	LANGUAGE plpgsql    
	AS $$
	BEGIN
		UPDATE validator
			SET uptime = uptime || 0
		WHERE power > 0;
	END;
	$$;`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
	CREATE OR REPLACE PROCEDURE increaseUptime(TEXT)
	LANGUAGE plpgsql    
	AS $$
	DECLARE
		len int;
	BEGIN
		SELECT
			array_length(uptime, 1) INTO len
		FROM validator
		WHERE consAddress = $1
		ORDER BY id DESC
		LIMIT 1;

		IF len > 0 THEN
			UPDATE validator
				SET uptime = uptime[1:array_length(uptime, 1)-1]
			WHERE id = (SELECT id FROM validator WHERE consAddress = $1 ORDER BY ID DESC LIMIT 1);
		END IF;	

		UPDATE validator
			SET uptime = uptime || 1
		WHERE id = (SELECT id FROM validator WHERE consAddress = $1 ORDER BY ID DESC LIMIT 1);

		SELECT
			array_length(uptime, 1) INTO len
		FROM validator
		WHERE consAddress = $1
		ORDER BY id DESC
		LIMIT 1;

		IF len > 100 THEN
			UPDATE validator
				SET uptime = uptime[2:101]
			WHERE id = (SELECT id FROM validator WHERE consAddress = $1 ORDER BY ID DESC LIMIT 1);
		END IF;
	END;
	$$;`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
	CREATE OR REPLACE PROCEDURE limitUptimeLen()
	LANGUAGE plpgsql    
	AS $$
	BEGIN
		UPDATE validator
			SET uptime = uptime[array_length(uptime, 1) - 99:array_length(uptime, 1)];
	END;
	$$;`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS transaction
	(
		id SERIAL PRIMARY KEY,
		height BIGINT,
		index INT,
		hash CHAR (64),
		code INTEGER,
		signers TEXT[],
		memo TEXT,
		tags JSONB,
		fee JSONB,
		gasWanted BIGINT,
		gasUsed BIGINT,
		msgs JSONB
	)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
	CREATE INDEX IF NOT EXISTS idx_tx_height_index
	ON transaction(height, index)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
	CREATE INDEX IF NOT EXISTS idx_tx_hash
	ON transaction
	USING HASH (hash)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
	CREATE INDEX IF NOT EXISTS idx_tx_gin_signers
	ON transaction
	USING GIN (signers)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
	CREATE INDEX IF NOT EXISTS idx_tx_gin_tags
	ON transaction
	USING GIN (tags)`)
	if err != nil {
		return err
	}

	return nil
}

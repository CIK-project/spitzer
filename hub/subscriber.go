package hub

import (
	"fmt"
	"encoding/hex"

	"database/sql"
	_ "github.com/lib/pq"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/CIK-project/spitzer/types"
	"github.com/CIK-project/spitzer/hub/model"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/cmd/gaia/app"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/codec"
)

// 코스모스-허브(가이아)용 섭스크라이버
type HubSubscriber struct {
	logger log.Logger
	db *sql.DB
	cdc *codec.Codec
	txDecoder sdk.TxDecoder
}

var _ types.Subscriber = &HubSubscriber{}

func NewHubSubscriber(logger log.Logger, config types.DBConfig) *HubSubscriber {
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
            config.User, config.Password, config.DBName)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		logger.Error(fmt.Sprintf("Error on open db: %s", err.Error()))
		panic(err)
	}

	cdc := app.MakeCodec()
	sdkConfig := sdk.GetConfig()
	sdkConfig.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	sdkConfig.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	sdkConfig.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	sdkConfig.Seal()

	hub := &HubSubscriber {
		logger: logger,
		db: db,
		cdc: cdc,
		txDecoder: auth.DefaultTxDecoder(cdc),
	}

	err = hub.initTables()
	if err != nil {
		logger.Error(fmt.Sprintf("Error on init tables: %s", err.Error()))
		panic(err)
	}
	
	return hub
}

func (hub *HubSubscriber) NextHeight() int64 {
	db := hub.db
	row := db.QueryRow(`
	SELECT
		height
	FROM
		block
	ORDER BY
		height DESC
	LIMIT 1		
	`)
	var lastHeight int64
	err := row.Scan(&lastHeight)
	if err != nil {
		hub.logger.Error(fmt.Sprintf("Error on query last height: %s", err.Error()))
		lastHeight = 0
	}

	return lastHeight + 1
}

func (hub *HubSubscriber) Commit(blockResult *types.BlockResult) error {
	dbTx, err := hub.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			hub.logger.Error(fmt.Sprintf("Panic on commit: %s", r))
			dbTx.Rollback()
		}
	}()

	block := blockResult.Block
	header := types.NewHeaderFrom(blockResult)
	r, err := model.SetHeader(dbTx, header)
	if err != nil {
		hub.logger.Error(fmt.Sprintf("Error on insert block: %s", err.Error()))
		return dbTx.Rollback()
	}
	hub.logger.Info("Insert block", "height", block.Height, "result", r)

	numTxs := block.NumTxs
	var i int64
	for i = 0; i < numTxs; i++ {
		tx, err := blockResult.GetTx(i)
		if err != nil {
			// TODO: error handling
			hub.logger.Error(fmt.Sprintf("Error get tx: %s", err.Error()))
			continue
		} else {
			var stdTx auth.StdTx
			appTx, sdkErr := hub.txDecoder(tx.Tx)
			// if err occured, just use default stdTx
			if sdkErr != nil {
				hub.logger.Error(fmt.Sprintf("Error on tx decoding: %s", sdkErr.Error()))
			} else {
				_stdTx, ok := appTx.(auth.StdTx)
				if ok == false {
					hub.logger.Error("Unkwon tx type")
				} else {
					stdTx = _stdTx
				}
			}
			stdTxJson, err := hub.cdc.MarshalJSON(stdTx)
			if err != nil {
				hub.logger.Error(fmt.Sprintf("Error on tx encoding json: %s", err.Error()))
				continue
			}

			r, err := dbTx.Exec(`
			INSERT INTO transaction(
				height,
				index,
				hash,
				code,
				stdTx
			)
			VALUES
				($1, $2, $3, $4, $5)
			`, block.Height,
			i,
			hex.EncodeToString(tx.Tx.Hash()),
			int32(tx.Result.Code),
			stdTxJson)
			if err != nil {
				hub.logger.Error(fmt.Sprintf("Error on insert transaction: %s", err.Error()))
				return dbTx.Rollback()
			}
			hub.logger.Info("Insert transaction", "height", block.Height, "index", i, "result", r)
		}
	}

	return dbTx.Commit()
}

func (hub *HubSubscriber) initTables() error {
	db := hub.db
	r, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS block
	(
		height BIGSERIAL PRIMARY KEY,
		hash CHAR (64) NOT NULL UNIQUE,
		prevHash CHAR (64) UNIQUE,
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
		proposer CHAR (40)
	)`)
	if err != nil {
		return err
	}
	hub.logger.Info("Create table", "result", r)

	db.Exec(`
	CREATE TABLE IF NOT EXISTS transaction
	(
		id SERIAL PRIMARY KEY,
		height BIGINT,
		index INT,
		hash CHAR (64),
		code INTEGER,
		stdTx json
	)`)

	return nil
}

func (hub *HubSubscriber) Stop() error {
	return hub.db.Close()
}

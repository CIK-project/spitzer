package hub

import (
	"fmt"
	"encoding/json"

	"database/sql"
	_ "github.com/lib/pq"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/CIK-project/spitzer/types"
	"github.com/CIK-project/spitzer/hub/model"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/cmd/gaia/app"
	"github.com/cosmos/cosmos-sdk/x/auth"
	stake "github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/codec"
	tmtypes "github.com/tendermint/tendermint/types"
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
	recentHeight, err := model.RecentHeight(hub.db)
	if err != nil {
		hub.logger.Error(fmt.Sprintf("Error on query recent height: %s", err.Error()))
		recentHeight = 0
	}

	return recentHeight + 1
}

func (hub *HubSubscriber) Genesis(genesis *tmtypes.GenesisDoc) error {
	genState := app.GenesisState{}
	json.Unmarshal(genesis.AppState, &genState)
	dbTx, err := hub.db.Begin()
	if err != nil {
		return err
	}

	for _, validator := range genState.StakingData.Validators {
		val := types.NewValidator(validator.OperatorAddr, validator.GetConsPubKey(), validator.GetPower().Int64())
		r, err := model.SetValidator(dbTx, val)
		if err != nil {
			return err
		}
		hub.logger.Info("Insert validator", "height", 0, "validator", val.Address, "result", r)
		
		desc := validator.Description
		r, err = model.SetValidatorDescription(dbTx, val.Address, desc.Moniker, desc.Identity, desc.Website, desc.Details)
		if err != nil {
			return err
		}
		hub.logger.Info("Insert validator description", "height", 0, "validator", val.Address, "moniker", validator.GetMoniker(), "result", r)
	}
	return dbTx.Commit()
}

func (hub *HubSubscriber) Commit(blockResult *types.BlockResult, catchUp bool) error {
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
		return err
	}
	hub.logger.Info("Insert block", "height", block.Height, "result", r)

	numTxs := block.NumTxs
	var i int64
	for i = 0; i < numTxs; i++ {
		txResult, err := blockResult.GetTx(i)
		if err != nil {
			// TODO: error handling, but never happen?
			hub.logger.Error(fmt.Sprintf("Error get tx: %s", err.Error()))
			continue
		} else {
			tx, err := types.NewTxFrom(&txResult, block.Height, i, hub.txDecoder, hub.cdc)
			if err != nil {
				return err
			}
			r, err := model.SetTx(dbTx, tx)
			if err != nil {
				return err
			}
			hub.logger.Info("Insert transaction", "height", block.Height, "index", i, "result", r)

			if txResult.Result.Code == 0 {
				tags := txResult.Result.GetTags()
				for _, tag := range tags {
					if string(tag.Key) == "action" && (string(tag.Value) == "create_validator" || string(tag.Value) == "edit_validator") {
						appTx, sdkErr := hub.txDecoder(txResult.Tx)
						if sdkErr == nil {
							stdTx, ok := appTx.(auth.StdTx)
							if ok {
								for _, msg := range stdTx.GetMsgs() {
									if msg, ok := msg.(stake.MsgCreateValidator); ok {
										val := types.NewValidator(msg.ValidatorAddr, msg.PubKey, 0)
										// first set validator
										_, err := model.SetValidator(dbTx, val)
										if err != nil {
											return err
										}
										desc := msg.Description
										r, err := model.SetValidatorDescription(dbTx, val.Address, desc.Moniker, desc.Identity, desc.Website, desc.Details)
										if err != nil {
											return err
										}
										hub.logger.Info("Set validator's description", "height", block.Height, "validator", val.Address, "moniker", desc.Moniker, "result", r)
									}

									if msg, ok := msg.(stake.MsgEditValidator); ok {
										val := &types.Validator {
											Address: msg.ValidatorAddr.String(),
										}

										desc := msg.Description
										r, err := model.SetValidatorDescription(dbTx, val.Address, desc.Moniker, desc.Identity, desc.Website, desc.Details)
										if err != nil {
											return err
										}
										hub.logger.Info("Set validator's description", "height", block.Height, "validator", val.Address, "moniker", desc.Moniker, "result", r)
									}

									// TODO: handle commision rate
								}
							}
						}
					}
				}
			}	
		}
	}

	endblock := blockResult.Results.EndBlock
	hub.logger.Info("Information of end block", "height", block.Height, "num valupdate", len(endblock.ValidatorUpdates))
	for _, valupdate := range endblock.ValidatorUpdates {
		pubkey, err := tmtypes.PB2TM.PubKey(valupdate.PubKey)
		if err != nil {
			return err
		}
		val := types.NewValidator(sdk.ValAddress{}, pubkey, valupdate.Power)
		r, err = model.SetValidatorPower(dbTx, val)
		if err != nil {
			return err
		}
		hub.logger.Info("Validator update", "height", block.Height, "validator_pub", val.PubKey, "power", val.Power, "result", r)
	}

	// calculate uptime only when catch up
	if catchUp {
		r, err = model.DecrementValsUptime(dbTx)
		if err != nil {
			return err
		}
		hub.logger.Info("Decrement all vals uptime", "height", block.Height, "result", r)

		numOfSigners := 0
		for _, precommit := range block.LastCommit.Precommits {
			if precommit != nil {
				val := types.Validator {
					ConsAddress: sdk.ConsAddress(precommit.ValidatorAddress).String(),
				}

				r, err = model.IncrementValUptime(dbTx, &val)
				if err != nil {
					return err
				}
				numOfSigners++
			}	
		}
		hub.logger.Info("Increment val's uptime", "height", block.Height, "num of signing vals", numOfSigners)

		_, err = model.LimitValsUptimeLen(dbTx)
		if err != nil {
			return err
		}
	}	

	return dbTx.Commit()
}

func (hub *HubSubscriber) initTables() error {
	return model.CreateTable(hub.db)
}

func (hub *HubSubscriber) Stop() error {
	return hub.db.Close()
}

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

type Validator struct {
	Address     string
	ConsAddress string
	PubKey      string
	Power       int64
	Moniker     string
	Identity    string
	Website     string
	Details     string
	Uptime      int
}

func NewValidator(valAddr sdk.ValAddress, pubKey crypto.PubKey, power int64) *Validator {
	return &Validator{
		Address:     valAddr.String(),
		ConsAddress: sdk.ConsAddress(pubKey.Address()).String(),
		PubKey:      sdk.MustBech32ifyValPub(pubKey),
		Power:       power,
	}
}

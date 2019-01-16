package types

import (
	tmtypes "github.com/tendermint/tendermint/types"
)

type Subscriber interface {
	NextHeight() int64
	Genesis(*tmtypes.GenesisDoc) error
	Commit(blockResult *BlockResult, catchUp bool) error
	Stop() error
}

type TestSubscriber struct {
	Height int64
}

var _ Subscriber = &TestSubscriber{}

func NewTestSubscriber(from int64) *TestSubscriber {
	return &TestSubscriber{
		Height: from,
	}
}

func (test *TestSubscriber) NextHeight() int64 {
	return test.Height + 1
}

func (test *TestSubscriber) Genesis(*tmtypes.GenesisDoc) error {
	return nil
}

func (test *TestSubscriber) Commit(*BlockResult, bool) error {
	test.Height++
	return nil
}

func (test *TestSubscriber) Stop() error {
	return nil
}

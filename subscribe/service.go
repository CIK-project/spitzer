package subscribe

import (
	"fmt"
	"time"

	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	rpc_client "github.com/tendermint/tendermint/rpc/client"

	"context"

	"github.com/CIK-project/spitzer/types"

	"github.com/tendermint/tendermint/libs/pubsub/query"
	tmtypes "github.com/tendermint/tendermint/types"
)

// 텐더민트를 섭스크라이브합니다.
// 텐더민트와 블록 높이가 동기화되어 있지 않으면 pulling을 합니다.
// 높이가 같아지면 catch up 상태가 되고 웹소켓으로 블록을 실시간으로 받아옵니다.
// TODO: Fail-over, adjusting exponential backoff...
type SubService struct {
	cmn.BaseService
	subscriber types.Subscriber
	caughtUp   bool

	client *rpc_client.HTTP
	wsCtx  context.Context
	wsOut  chan interface{}
}

func NewSubService(logger log.Logger, subscriber types.Subscriber, remote, wsEndpoint string) *SubService {
	ss := &SubService{}
	ss.subscriber = subscriber
	ss.client = rpc_client.NewHTTP(remote, wsEndpoint)
	ss.client.SetLogger(logger)
	ss.client.WSEvents.SetLogger(logger)

	ss.BaseService = *cmn.NewBaseService(logger, "Subscription service", ss)
	return ss
}

func (ss *SubService) OnStart() error {
	err := ss.BaseService.OnStart()
	if err != nil {
		return err
	}

	err = ss.client.Start()
	if err != nil {
		return err
	}

	ss.wsCtx = context.Background()
	ss.wsOut = make(chan interface{})
	q, err := query.New("tm.event='NewBlock'")
	if err != nil {
		panic(err)
	}
	err = ss.client.Subscribe(ss.wsCtx, "new block", q, ss.wsOut)
	if err != nil {
		panic(err)
	}

	go ss.routine()

	return nil
}

func (ss *SubService) routine() {
	for {
		select {
		case i, ok := <-ss.wsOut:
			if ok {
				ss.wsRoutine(i)
			}
		default:
			ss.pullRoutine()
		}
	}
}

func (ss *SubService) pullRoutine() {
	if ss.caughtUp {
		time.Sleep(100 * time.Microsecond)
		return
	}

	height := ss.subscriber.NextHeight()
	if height == 1 {
		genesis, err := ss.client.Genesis()
		if err != nil {
			ss.Logger.Error(fmt.Sprintf("Error on genesis rpc: %s", err.Error()), "height", 0)
			return
		}

		err = ss.subscriber.Genesis(genesis.Genesis)
		if err != nil {
			ss.Logger.Error(fmt.Sprintf("Error on genesis: %s", err.Error()), "height", 0)
			return
		}
	}

	block, err := ss.client.Block(&height)
	if err != nil {
		ss.Logger.Error(fmt.Sprintf("Error on pulling: %s", err.Error()), "height", ss.subscriber.NextHeight())
		time.Sleep(1 * time.Second)
		return
	}

	result, err := ss.client.BlockResults(&height)
	if err != nil {
		ss.Logger.Error(fmt.Sprintf("Error on pulling: %s", err.Error()), "height", ss.subscriber.NextHeight())
		time.Sleep(1 * time.Second)
		return
	}

	blockResult := &types.BlockResult{
		Block:   block.Block,
		Results: result.Results,
	}

	ss.Logger.Info("Commit block from pulling", "height", block.Block.Height)
	err = ss.subscriber.Commit(blockResult, false)
	if err != nil {
		ss.Logger.Error(fmt.Sprintf("Fail to commit block: %s", err.Error()))
	}
}

func (ss *SubService) wsRoutine(i interface{}) {
	switch block := i.(type) {
	case tmtypes.EventDataNewBlock:

		if block.Block.Height == ss.subscriber.NextHeight() {
			if ss.caughtUp == false {
				ss.caughtUp = true
				ss.Logger.Info("Catch up!", "height", block.Block.Height)
			}
		} else {
			ss.Logger.Info("Get block from ws, but not expected", "expected", ss.subscriber.NextHeight(), "but", block.Block.Height)
			if ss.caughtUp {
				ss.caughtUp = false
				ss.Logger.Info("Lost catching up", "expected", ss.subscriber.NextHeight(), "but", block.Block.Height)
			}
			return
		}

		if block.Block.Height == 1 {
			genesis, err := ss.client.Genesis()
			if err != nil {
				ss.Logger.Error(fmt.Sprintf("Error on genesis rpc: %s", err.Error()), "height", 0)
				return
			}
	
			err = ss.subscriber.Genesis(genesis.Genesis)
			if err != nil {
				ss.Logger.Error(fmt.Sprintf("Error on genesis: %s", err.Error()), "height", 0)
				return
			}
		}

		result, err := ss.client.BlockResults(&block.Block.Height)
		if err != nil {
			ss.Logger.Error(err.Error())
		} else {
			blockResult := &types.BlockResult{
				Block:   block.Block,
				Results: result.Results,
			}

			ss.Logger.Info("Commit block from ws", "height", block.Block.Height)
			err = ss.subscriber.Commit(blockResult, true)
			if err != nil {
				ss.Logger.Error(fmt.Sprintf("Fail to commit block: %s", err.Error()))
			}
		}
	default:
		ss.Logger.Error("Unknown new block type")
	}
}

func (ss *SubService) OnStop() {
	ss.BaseService.OnStop()
	ss.client.OnStop()

	if err := ss.subscriber.Stop(); err != nil {
		ss.Logger.Error(fmt.Sprintf("Error on stop subscriber: %s", err.Error()))
	}
}

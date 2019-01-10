package gql

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"

	"github.com/CIK-project/spitzer/types"
	"github.com/tendermint/tendermint/libs/log"
	cmn "github.com/tendermint/tendermint/libs/common"
)

type GqlService struct {
	cmn.BaseService

	server *Server
}

func NewGqlService(logger log.Logger, config types.DBConfig) *GqlService {
	gs := &GqlService{}

	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
            config.User, config.Password, config.DBName)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		logger.Error(fmt.Sprintf("Error on open db: %s", err.Error()))
		panic(err)
	}
	gs.server = NewServer(logger, db)

	gs.BaseService = *cmn.NewBaseService(logger, "Graphql service", gs)
	return gs
}

func (gs *GqlService) OnStart() error {
	return gs.server.Start()
}

func (gs *GqlService) OnStop() {
	gs.BaseService.OnStop()
}

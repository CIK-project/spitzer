package schema

import (
	"fmt"
	"errors"
	"database/sql"

	"github.com/graphql-go/graphql"
	"github.com/CIK-project/spitzer/types"

	"github.com/CIK-project/spitzer/hub/model"
)

var HeaderObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "Header",
	Fields: graphql.Fields{
		"height": &graphql.Field{
			Type: graphql.Int,
		},
		"hash": &graphql.Field{
			Type: graphql.String,
		},
		"prevHash": &graphql.Field{
			Type: graphql.String,
		},
		"time": &graphql.Field{
			Type: graphql.DateTime,
		},
		"numTxs": &graphql.Field{
			Type: graphql.Int,
		},
		"totalTxs": &graphql.Field{
			Type: graphql.Int,
		},
		"lastCommitHash": &graphql.Field{
			Type: graphql.String,
		},
		"dataHash": &graphql.Field{
			Type: graphql.String,
		},
		"validatorHash": &graphql.Field{
			Type: graphql.String,
		},
		"nextValidatorHash": &graphql.Field{
			Type: graphql.String,
		},
		"consensusHash": &graphql.Field{
			Type: graphql.String,
		},
		"appHash": &graphql.Field{
			Type: graphql.String,
		},
		"lastResultHash": &graphql.Field{
			Type: graphql.String,
		},
		"evidenceHash": &graphql.Field{
			Type: graphql.String,
		},
		"proposer": &graphql.Field{
			Type: ValidatorObject,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				db, ok := p.Context.Value(ContextKeyDB).(*sql.DB)
				if !ok {
					return nil, errors.New("Invalid context")
				}
				header, ok := p.Source.(*types.Header)
				if !ok {
					return nil, errors.New("Invalid header type")
				}
				val, err := model.GetValidatorByCons(db, header.Proposer)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("Can't get validator: %s, consAddress: %s", err.Error(), header.Proposer))
				}

				return val, nil
			},
		},
		"tags": &graphql.Field{
			Type: graphql.String,
		},
	},
})

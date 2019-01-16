package schema

import (
	"database/sql"
	"errors"
	"github.com/graphql-go/graphql"
	"fmt"

	"github.com/CIK-project/spitzer/hub/model"
)

var QueryObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"header": &graphql.Field{
			Type: HeaderObject,
			Description: "Get specific header at given height",
			Args: graphql.FieldConfigArgument{
				"height": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				height, ok := p.Args["height"].(int)
				if !ok {
					return nil, errors.New("Invalid height type")
				}
				db, ok := p.Context.Value(ContextKeyDB).(*sql.DB)
				if !ok {
					return nil, errors.New("Invalid context")
				}
				
				header, err := model.GetHeader(db, int64(height))
				if err != nil {
					return nil, errors.New(fmt.Sprintf("Can't get header: %s", err.Error()))
				}
				return header, nil
			},
		},
		"headers": &graphql.Field{
			Type: graphql.NewList(HeaderObject),
			Description: "Get headers by descending height",
			Args: graphql.FieldConfigArgument{
				"page": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"limit": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				page, ok := p.Args["page"].(int)
				if !ok {
					return nil, errors.New("Invalid page type")
				}
				limit, ok := p.Args["limit"].(int)
				if !ok {
					return nil, errors.New("Invalid limit type")
				}
				db, ok := p.Context.Value(ContextKeyDB).(*sql.DB)
				if !ok {
					return nil, errors.New("Invalid context")
				}
				
				headers, err := model.GetHeaders(db, page, limit)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("Can't get headers: %s", err.Error()))
				}
				return headers, nil
			},
		},
		"txs": &graphql.Field{
			Type: graphql.NewList(TxObject),
			Description: `Get transactions by decending height and index.
You can use specific height as param.
Signers parameter act as OR, example usage is {cosmos1xf6upkfzgl8f3pkpa0l2d76pex33ur277x5vwj, cosmos1l3agyg568dsk2r9v7f35wge3yr85ampz6cy720}.
Tags paramater act as And, example usage is [{"action":"delegate"}].`,
			Args: graphql.FieldConfigArgument {
				"page": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"limit": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"height": &graphql.ArgumentConfig{
					Type: graphql.Int,
					DefaultValue: 0,
				},
				"signers": &graphql.ArgumentConfig{
					Type: graphql.String,
					DefaultValue: "{}",
				},
				"tags": &graphql.ArgumentConfig{
					Type: graphql.String,
					DefaultValue: "[]",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				page, ok := p.Args["page"].(int)
				if !ok {
					return nil, errors.New("Invalid page type")
				}
				limit, ok := p.Args["limit"].(int)
				if !ok {
					return nil, errors.New("Invalid limit type")
				}
				height, ok := p.Args["height"].(int)
				if !ok {
					return nil, errors.New("Invalid height type")
				}
				signers, ok := p.Args["signers"].(string)
				if !ok {
					return nil, errors.New("Invalid signers type")
				}
				tags, ok := p.Args["tags"].(string)
				if !ok {
					return nil, errors.New("Invalid tags type")
				}
				db, ok := p.Context.Value(ContextKeyDB).(*sql.DB)
				if !ok {
					return nil, errors.New("Invalid context")
				}
				
				txs, err := model.GetTxs(db, page, limit, int64(height), signers, tags)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("Can't get txs: %s", err.Error()))
				}
				return txs, nil
			},
		},
		"validator": &graphql.Field{
			Type: ValidatorObject,
			Description: "Get validators by val address",
			Args: graphql.FieldConfigArgument{
				"address": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				address, ok := p.Args["address"].(string)
				if !ok {
					return nil, errors.New("Invalid address type")
				}
				db, ok := p.Context.Value(ContextKeyDB).(*sql.DB)
				if !ok {
					return nil, errors.New("Invalid context")
				}
				
				val, err := model.GetValidator(db, address)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("Can't get validator: %s", err.Error()))
				}
				return val, nil
			},
		},
		"validators": &graphql.Field{
			Type: graphql.NewList(ValidatorObject),
			Description: "Get validators by decending their power.",
			Args: graphql.FieldConfigArgument{
				"page": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"limit": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				page, ok := p.Args["page"].(int)
				if !ok {
					return nil, errors.New("Invalid page type")
				}
				limit, ok := p.Args["limit"].(int)
				if !ok {
					return nil, errors.New("Invalid limit type")
				}
				db, ok := p.Context.Value(ContextKeyDB).(*sql.DB)
				if !ok {
					return nil, errors.New("Invalid context")
				}
				
				vals, err := model.GetValidators(db, page, limit)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("Can't get validators: %s", err.Error()))
				}
				return vals, nil
			},
		},
		"numOfValidators": &graphql.Field{
			Type: graphql.Int,
			Description: "Get number of validators.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				db, ok := p.Context.Value(ContextKeyDB).(*sql.DB)
				if !ok {
					return nil, errors.New("Invalid context")
				}
				
				num, err := model.GetNumOfValidators(db)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("Can't get num of validators: %s", err.Error()))
				}
				return num, nil
			},
		},
	},
})
package schema

import (
	"github.com/graphql-go/graphql"
)

var TxObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "Transaction",
	Fields: graphql.Fields{
		"height": &graphql.Field{
			Type: graphql.Int,
		},
		"index": &graphql.Field{
			Type: graphql.Int,
		},
		"hash": &graphql.Field{
			Type: graphql.String,
		},
		"code": &graphql.Field{
			Type: graphql.Int,
		},
		"signers": &graphql.Field{
			Type: graphql.NewList(graphql.String),
		},
		"memo": &graphql.Field{
			Type: graphql.String,
		},
		"tags": &graphql.Field{
			Type: graphql.String,
		},
		"fee": &graphql.Field{
			Type: graphql.String,
		},
		"gasWanted": &graphql.Field{
			Type: graphql.Int,
		},
		"gasUsed": &graphql.Field{
			Type: graphql.Int,
		},
		"msgs": &graphql.Field{
			Type: graphql.String,
		},
	},
})

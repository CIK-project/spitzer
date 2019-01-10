package types

import (
	"github.com/graphql-go/graphql"
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
			Type: graphql.String,
		},
	},
})

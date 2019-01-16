package schema

import (
	"github.com/graphql-go/graphql"
)

var ValidatorObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "Validator",
	Fields: graphql.Fields{
		"address": &graphql.Field{
			Type: graphql.String,
		},
		"consAddress": &graphql.Field{
			Type: graphql.String,
		},
		"pubKey": &graphql.Field{
			Type: graphql.String,
		},
		"power": &graphql.Field{
			Type: graphql.Int,
		},
		"moniker": &graphql.Field{
			Type: graphql.String,
		},
		"identity": &graphql.Field{
			Type: graphql.String,
		},
		"website": &graphql.Field{
			Type: graphql.String,
		},
		"details": &graphql.Field{
			Type: graphql.String,
		},
		"uptime": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

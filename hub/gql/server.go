package gql

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"errors"
	"fmt"

	"github.com/CIK-project/spitzer/hub/model"
	gqltypes "github.com/CIK-project/spitzer/hub/gql/types"
	"github.com/graphql-go/graphql"
	"github.com/tendermint/tendermint/libs/log"
)

type Server struct {
	logger log.Logger
	DB     *sql.DB
	schema *graphql.Schema
}

func NewServer(logger log.Logger, db *sql.DB) *Server {
	return &Server{
		logger: logger,
		DB:     db,
	}
}

func (s *Server) Start() error {
	queryObject := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"header": &graphql.Field{
				Type: gqltypes.HeaderObject,
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
					
					header, err := model.GetHeader(s.DB, int64(height))
					if err != nil {
						s.logger.Error(fmt.Sprintf("Error on get header: %s", err.Error()))
						return nil, errors.New("Can't get header")
					}
					return header, nil
				},
			},
			"headers": &graphql.Field{
				Type: graphql.NewList(gqltypes.HeaderObject),
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
					
					headers, err := model.GetHeaders(s.DB, page, limit)
					if err != nil {
						s.logger.Error(fmt.Sprintf("Error on get headers: %s", err.Error()))
						return nil, errors.New("Can't get headers")
					}
					return headers, nil
				},
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryObject,
	})
	if err != nil {
		return err
	}
	s.schema = &schema

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := graphql.Do(graphql.Params{
			Schema:        *s.schema,
			RequestString: r.URL.Query().Get("query"),
		})
		json.NewEncoder(w).Encode(result)
	})

	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			panic(err)
		}
	}()
	return nil
}

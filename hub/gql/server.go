package gql

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	schematypes "github.com/CIK-project/spitzer/hub/gql/schema"
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
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: schematypes.QueryObject,
	})
	if err != nil {
		return err
	}
	s.schema = &schema

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*r).Method == "OPTIONS" {
			return
		}

		result := graphql.Do(graphql.Params{
			Schema:        *s.schema,
			Context:       context.WithValue(context.Background(), schematypes.ContextKeyDB, s.DB),
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

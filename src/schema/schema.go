package schema

import (
	"context"
	"errors"
	"fmt"

	"github.com/PaesslerAG/gval"
	"github.com/graphql-go/graphql"
	goke "github.com/remorses/goke/src"
	"github.com/remorses/goke/src/fakedata"
	"github.com/remorses/goke/src/firestore"
	"github.com/remorses/goke/src/mongodb"
	"github.com/remorses/goke/src/types"
	tools "github.com/remorses/graphql-go-tools"
)

// MakeGokeSchema generates the schema
func MakeGokeSchema(config goke.Config) (graphql.Schema, error) {
	config, err := InitializeConfig(config)
	if err != nil {
		return graphql.Schema{}, err
	}
	schema, err := generateSchema(config)
	if err != nil {
		return schema, err
	}
	return schema, nil
}

func InitializeConfig(config goke.Config) (goke.Config, error) {
	err := config.Init()
	if err != nil {
		return config, err
	}
	// database functions
	if config.Mongodb.Uri != "" {
		config.DatabaseFunctions = &mongodb.MongodbDatabaseFunctions{Config: config}
	} else if config.Firestore.ProjectID != "" {
		config.DatabaseFunctions = &firestore.FirestoreDatabaseFunctions{Config: config}
	}

	// by default use local fake data
	if config.DatabaseFunctions == nil {
		println("using local fake data source")
		config.DatabaseFunctions = &fakedata.FakeDatabaseFunctions{Config: config}
	}
	return config, nil
}

func makeSchemaConfig(config goke.Config) (graphql.SchemaConfig, error) {
	resolvers := map[string]tools.Resolver{
		types.ObjectID.Name(): &tools.ScalarResolver{
			Serialize:    types.ObjectID.Serialize,
			ParseLiteral: types.ObjectID.ParseLiteral,
			ParseValue:   types.ObjectID.ParseValue,
		},
	}
	// IsTypeOf
	for name, typeConf := range config.Types {
		if typeConf.IsTypeOf == "" {
			continue
		}
		eval, err := gval.Full().NewEvaluable(typeConf.IsTypeOf)
		if err != nil {
			return graphql.SchemaConfig{}, errors.New(
				"got an error parsing isTypeOf expression " + typeConf.IsTypeOf)
		}
		resolvers[name] = &tools.ObjectResolver{
			IsTypeOf: func(p graphql.IsTypeOfParams) bool {
				res, err := eval(context.Background(), goke.Map{
					"x":        p.Value, // TODO add more variables to isTypeOf expressions
					"document": p.Value,
				})
				if err != nil {
					fmt.Println("got an error evaluating expression " + typeConf.IsTypeOf)
					return false
				}
				if res == true {
					return true
				}
				return false
			},
		}
	}

	baseSchemaConfig, err := tools.MakeSchemaConfig(
		tools.ExecutableSchema{
			TypeDefs:  []string{config.Schema},
			Resolvers: resolvers,
		},
	)
	return baseSchemaConfig, err
}

func generateSchema(config goke.Config) (graphql.Schema, error) {
	baseSchemaConfig, err := makeSchemaConfig(config)
	if err != nil {
		return graphql.Schema{}, err
	}
	// add fields
	query, err := makeQuery(config, baseSchemaConfig)
	if err != nil {
		return graphql.Schema{}, err
	}
	mutation, err := makeMutation(config, baseSchemaConfig)
	if err != nil {
		return graphql.Schema{}, err
	}
	err = addRelationsFields(config, baseSchemaConfig)
	if err != nil {
		return graphql.Schema{}, err
	}

	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Types:      baseSchemaConfig.Types,
			Extensions: baseSchemaConfig.Extensions,
			Query:      query,
			Mutation:   mutation,
		},
	)
	if err != nil {
		return graphql.Schema{}, err
	}
	return schema, nil
}

func unexposedType(typeConf *goke.TypeConfig) bool {
	return typeConf == nil || (typeConf.Exposed != nil && !*typeConf.Exposed)
}

func findType(a []graphql.Type, name string) graphql.Type {
	for _, t := range a {
		if t.Name() == name {
			return t
		}
	}
	return nil
}

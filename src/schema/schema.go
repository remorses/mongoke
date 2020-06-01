package schema

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/PaesslerAG/gval"
	"github.com/graphql-go/graphql"
	tools "github.com/remorses/graphql-go-tools"
	mongoke "github.com/remorses/mongoke/src"
	"github.com/remorses/mongoke/src/fakedata"
	"github.com/remorses/mongoke/src/firestore"
	"github.com/remorses/mongoke/src/mongodb"
	"github.com/remorses/mongoke/src/types"
)

// MakeMongokeSchema generates the schema
func MakeMongokeSchema(config mongoke.Config) (graphql.Schema, error) {

	if config.Cache == nil {
		config.Cache = make(mongoke.Map)
	}

	if config.Schema == "" && config.SchemaPath != "" {
		data, e := ioutil.ReadFile(config.SchemaPath)
		if e != nil {
			return graphql.Schema{}, e
		}
		config.Schema = string(data)
	}

	if config.Schema == "" && config.SchemaUrl != "" {
		data, e := mongoke.DownloadFile(config.SchemaUrl)
		if e != nil {
			return graphql.Schema{}, e
		}
		config.Schema = string(data)
	}

	if config.Schema == "" {
		return graphql.Schema{}, errors.New("missing required schema")
	}

	if config.Mongodb.Uri != "" {
		config.DatabaseFunctions = &mongodb.MongodbDatabaseFunctions{Config: config}
	} else if config.Firestore.ProjectID != "" {
		config.DatabaseFunctions = &firestore.FirestoreDatabaseFunctions{Config: config}
	}

	if config.DatabaseFunctions == nil { // by default use local fake data
		println("using local fake data source")
		config.DatabaseFunctions = &fakedata.FakeDatabaseFunctions{Config: config}
	}

	schema, err := generateSchema(config)
	if err != nil {
		return schema, err
	}
	return schema, nil
}

func makeSchemaConfig(config mongoke.Config) (graphql.SchemaConfig, error) {
	resolvers := map[string]tools.Resolver{
		types.ObjectID.Name(): &tools.ScalarResolver{
			Serialize:    types.ObjectID.Serialize,
			ParseLiteral: types.ObjectID.ParseLiteral,
			ParseValue:   types.ObjectID.ParseValue,
		},
	}
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
				res, err := eval(context.Background(), mongoke.Map{
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

func generateSchema(config mongoke.Config) (graphql.Schema, error) {
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

func unexposedType(typeConf *mongoke.TypeConfig) bool {
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

package mongoke

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type Mongoke struct {
	databaseFunctions DatabaseInterface
	typeDefs          string
	databaseUri       string
	typeMap           map[string]graphql.Type
	Config            Config
}

// MakeMongokeSchema generates the schema
func MakeMongokeSchema(config Config, databaseFunctions DatabaseInterface) (graphql.Schema, error) {
	if databaseFunctions == nil {
		databaseFunctions = MongodbDatabaseFunctions{}
	}
	if config.Schema == "" && config.SchemaPath != "" {
		data, e := ioutil.ReadFile(config.SchemaPath)
		if e != nil {
			return graphql.Schema{}, e
		}
		config.Schema = string(data)
	}
	mongoke := Mongoke{
		Config:            config,
		typeDefs:          config.Schema,
		databaseFunctions: databaseFunctions,
		typeMap:           make(map[string]graphql.Type),
		databaseUri:       config.DatabaseUri,
	}
	schema, err := mongoke.generateSchema()
	if err != nil {
		return schema, err
	}
	return schema, nil
}

// MakeMongokeHandler creates an http handler
func MakeMongokeHandler(config Config, databaseFunctions DatabaseInterface) (http.Handler, error) {
	schema, err := MakeMongokeSchema(config, databaseFunctions)
	if err != nil {
		return nil, err
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
		RootObjectFn: func(ctx context.Context, r *http.Request) map[string]interface{} {
			// TODO add user jwt data here, resolver can return an error if user not authenticated
			return make(map[string]interface{})
		},
	})
	return h, nil
}

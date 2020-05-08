package mongoke

import (
	"context"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type Mongoke struct {
	typeDefs   string
	mongoDbUri string
	typeMap    map[string]graphql.Type
}

// MakeMongokeSchema generates the schema
func MakeMongokeSchema(config Config) (graphql.Schema, error) {
	mongoke := Mongoke{
		typeDefs:   config.schemaString,
		typeMap:    make(map[string]graphql.Type),
		mongoDbUri: config.mongoDbUri,
	}
	schema, err := mongoke.generateSchema()
	if err != nil {
		return schema, err
	}
	return schema, nil
}

// TODO read the yaml and parse into Config

func main(config Config) {
	schema, err := MakeMongokeSchema(config)

	if err != nil {
		panic(err)
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

	http.Handle("/", h)
	http.ListenAndServe("localhost:8080", nil)
}

package mongoke

import (
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type Config struct {
	schemaString string
	mongoDbUri   string
}

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

func main(config Config) {
	schema, err := MakeMongokeSchema(config)

	if err != nil {
		panic(err)
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/", h)
	http.ListenAndServe("localhost:8080", nil)
}

package mongoke

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PaesslerAG/gval"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	tools "github.com/remorses/graphql-go-tools"
)

type Mongoke struct {
	databaseFunctions  DatabaseInterface
	typeDefs           string
	databaseUri        string
	indexableTypeNames []string
	typeMap            map[string]graphql.Type
	Config             Config
	schemaConfig       graphql.SchemaConfig
}

// MakeMongokeSchema generates the schema
func MakeMongokeSchema(config Config) (graphql.Schema, error) {
	if config.databaseFunctions == nil {
		config.databaseFunctions = MongodbDatabaseFunctions{}
	}
	// TODO validate config here

	if config.Schema == "" && config.SchemaPath != "" {
		data, e := ioutil.ReadFile(config.SchemaPath)
		if e != nil {
			return graphql.Schema{}, e
		}
		config.Schema = string(data)
	}
	schemaConfig, err := makeSchemaConfig(config)
	if err != nil {
		return graphql.Schema{}, err
	}
	mongoke := Mongoke{
		Config:            config,
		typeDefs:          config.Schema,
		databaseFunctions: config.databaseFunctions,
		typeMap:           make(map[string]graphql.Type),
		databaseUri:       config.DatabaseUri,
		schemaConfig:      schemaConfig,
	}
	schema, err := mongoke.generateSchema()
	if err != nil {
		return schema, err
	}
	return schema, nil
}

func makeSchemaConfig(config Config) (graphql.SchemaConfig, error) {
	resolvers := map[string]tools.Resolver{
		objectID.Name(): &tools.ScalarResolver{
			Serialize:    objectID.Serialize,
			ParseLiteral: objectID.ParseLiteral,
			ParseValue:   objectID.ParseValue,
		},
	}
	for name, typeConf := range config.Types {
		if typeConf.IsTypeOf == "" {
			continue
		}
		eval, err := gval.Full().NewEvaluable(typeConf.IsTypeOf)
		if err != nil {
			return graphql.SchemaConfig{}, errors.New("got an error parsing isTypeOf expression " + typeConf.IsTypeOf)
		}
		resolvers[name] = &tools.ObjectResolver{
			IsTypeOf: func(p graphql.IsTypeOfParams) bool {
				res, err := eval(context.Background(), Map{
					"x":        p.Value,
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

// MakeMongokeHandler creates an http handler
func MakeMongokeHandler(config Config) (http.Handler, error) {
	schema, err := MakeMongokeSchema(config)
	if err != nil {
		return nil, err
	}

	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
		// GraphiQL:   true,
		Playground: true,
		RootObjectFn: func(ctx context.Context, r *http.Request) map[string]interface{} {
			rootValue := Map{
				"request": r,
			}

			tknStr := r.Header.Get("Authorization")
			parts := strings.Split(tknStr, "Bearer")
			tknStr = reverseStrings(parts)[0]
			tknStr = strings.TrimSpace(tknStr)

			if tknStr == "" {
				return rootValue
			}

			claims, err := extractClaims(config.JwtConfig, tknStr)

			if err != nil {
				fmt.Println("error in handler", err)
				return rootValue
			}

			rootValue["jwt"] = claims

			return rootValue
		},
	})

	return h, nil
}

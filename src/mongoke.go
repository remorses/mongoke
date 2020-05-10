package mongoke

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"

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
	if config.Schema == "" && config.SchemaPath != "" {
		data, e := ioutil.ReadFile(config.SchemaPath)
		if e != nil {
			return graphql.Schema{}, e
		}
		config.Schema = string(data)
	}
	schemaConfig, err := makeSchemaConfig(config.Schema)
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

func makeSchemaConfig(typeDefs string) (graphql.SchemaConfig, error) {
	baseSchemaConfig, err := tools.MakeSchemaConfig(
		tools.ExecutableSchema{
			TypeDefs: []string{typeDefs},
			Resolvers: map[string]tools.Resolver{
				objectID.Name(): &tools.ScalarResolver{
					Serialize:    objectID.Serialize,
					ParseLiteral: objectID.ParseLiteral,
					ParseValue:   objectID.ParseValue,
				},
			},
		},
	)
	return baseSchemaConfig, err
}

// MakeMongokeHandler creates an http handler
func MakeMongokeHandler(config Config) (http.HandlerFunc, error) {
	schema, err := MakeMongokeSchema(config)
	if err != nil {
		return nil, err
	}

	h := func(w http.ResponseWriter, r *http.Request) {

		// parse http.Request into handler.RequestOptions
		opts := handler.NewRequestOptions(r)

		tknStr := r.Header.Get("Authorization")
		parts := strings.Split(tknStr, "Bearer")
		tknStr = reverseStrings(parts)[0]
		tknStr = strings.TrimSpace(tknStr)
		claims, err := extractClaims(tknStr, "secret") // TODO take secret from config or url

		rootValue := Map{
			"response": w,
			"request":  r,
			"jwt":      claims,
		}

		// execute graphql query
		// here, we passed in Query, Variables and OperationName extracted from http.Request
		params := graphql.Params{
			Schema:         schema,
			RequestString:  opts.Query,
			VariableValues: opts.Variables,
			OperationName:  opts.OperationName,
			RootObject:     rootValue,
		}
		result := graphql.Do(params)

		// one way to render JSON without use of external libraries
		// alternatively, you can use libraries like `unrolled/render`
		js, err := json.Marshal(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}

	return h, nil
}

func extractClaims(tokenStr string, secret string) (jwt.MapClaims, error) {
	hmacSecret := []byte(secret)
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// TODO check token ALG method is same etc
		return hmacSecret, nil
	})

	if err != nil {
		return jwt.MapClaims{}, err
	}

	return claims, nil

}

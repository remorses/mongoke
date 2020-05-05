package mongoke

import (
	"encoding/json"
	"fmt"
	"log"

	tools "github.com/bhoriuchi/graphql-go-tools"
	"github.com/graphql-go/graphql"
)

type Config struct {
	schemaString string
}

func main(config Config) {
	schema, err := tools.MakeExecutableSchema(tools.ExecutableSchema{
		TypeDefs: config.schemaString,
		Resolvers: tools.ResolverMap{
			"Query": &tools.ObjectResolver{
				Fields: tools.FieldResolveMap{
					"foo": func(p graphql.ResolveParams) (interface{}, error) {
						// lookup data
						return nil, nil
					},
				},
			},
		},
	})

	// for _, field := range schema.QueryType().Fields() {
	// 	name := field.Name
	// 	println(name)
	// }

	if err != nil {
		log.Fatalf("Failed to build schema, error: %v", err)
	}
	println(schema.TypeMap())
	params := graphql.Params{
		Schema: schema,
		RequestString: `
	  query GetFoo {
		foo(id: "5cffbf1ccecefcfff659cea8") {
		  description
		}
	  }`,
	}

	r := graphql.Do(params)
	if r.HasErrors() {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON)
}

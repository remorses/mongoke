package benchmark

import (
	"context"
	"testing"

	_ "net/http/pprof"

	"github.com/graphql-go/graphql"
	mongoke "github.com/remorses/mongoke/src"
	"github.com/remorses/mongoke/src/mock"
	"github.com/remorses/mongoke/src/schema"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func BenchmarkQuery(b *testing.B) {

	db := &mock.DatabaseInterfaceMock{
		FindManyFunc: func(ctx context.Context, p mongoke.FindManyParams) ([]map[string]interface{}, error) {
			elem := mongoke.Map{
				"_id":  primitive.NewObjectID(),
				"name": "1",
				"nested": mongoke.Map{
					"x":     1,
					"field": "ciao",
				},
				"age": 10,
			}
			return append(
				[]mongoke.Map{},
				elem,
				elem,
				elem,
				elem,
				elem,
				elem,
				elem,
				elem,
				elem,
			), nil
		},
	}

	s, _ := schema.MakeMongokeSchema(mongoke.Config{
		Schema: `
		scalar ObjectId
		interface Named {
			name: String
		}

		type User implements Named {
			_id: ObjectId!
			name: String
			age: Int!
			nested: Obj
		}

		type Obj {
			field: String
			x: Int
		}
		`,
		DatabaseFunctions: db,
		Types: map[string]*mongoke.TypeConfig{
			"User": {Collection: "users"},
		},
	})

	b.Run("main", func(b *testing.B) {

		for i := 0; i < b.N; i++ {
			res := graphql.Do(graphql.Params{
				Schema: s,
				RequestString: `
			{
				User(where: {name: {neq: "xxx"}}) {
					name
					age
					_id
					nested {
						field
						x
					}
				}
			}
			`,
			})
			if len(res.Errors) != 0 {
				b.Fatal(res.Errors[0])
			}
		}
	})

}

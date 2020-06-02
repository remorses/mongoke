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
			return []mongoke.Map{
				{
					"_id":  primitive.NewObjectID(),
					"name": "1",
					"age":  10,
				},
			}, nil
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
				User {
					name
					age
					_id
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

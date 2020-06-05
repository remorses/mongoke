package schema_test

import (
	"context"
	"errors"
	"testing"

	"github.com/graphql-go/graphql"
	mongoke "github.com/remorses/mongoke/src"
	"github.com/remorses/mongoke/src/fakedata"
	"github.com/remorses/mongoke/src/mock"
	mongoke_schema "github.com/remorses/mongoke/src/schema"
	"github.com/remorses/mongoke/src/testutil"
)

func TestMutationWithEmptyFakeDatabase(t *testing.T) {
	db := &fakedata.FakeDatabaseFunctions{}
	schema, _ := mongoke_schema.MakeMongokeSchema(mongoke.Config{
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

	testutil.NewTestGroup(t, testutil.NewTestGroupParams{
		Collection:    "users",
		Database:      db,
		Documents:     []mongoke.Map{},
		DefaultSchema: schema,
		Tests: []testutil.TestCase{
			{
				Name:     "update with set",
				Schema:   schema,
				Expected: mongoke.Map{"updateUser": mongoke.Map{"returning": nil, "affectedCount": 0}},
				Query: `
				mutation {
					updateUser(set: {name: "xxx"}) {
						affectedCount
						returning {
							name
							age
							_id
						}
					}
				}
				`,
			},
			{
				Name:     "update with set and where",
				Schema:   schema,
				Expected: mongoke.Map{"updateUser": mongoke.Map{"returning": nil, "affectedCount": 0}},
				Query: `
				mutation {
					updateUser(where: {name: {eq: "zzz"}}, set: {name: "xxx"}) {
						affectedCount
						returning {
							name
							age
							_id
						}
					}
				}
				`,
			},
			{
				Name:     "updateNodes with set",
				Schema:   schema,
				Expected: mongoke.Map{"updateUserNodes": mongoke.Map{"returning": []mongoke.Map{}, "affectedCount": 0}},
				Query: `
				mutation {
					updateUserNodes(set: {name: "xxx"}) {
						affectedCount
						returning {
							name
							age
							_id
						}
					}
				}
				`,
			},
			{
				Name:     "insert",
				Schema:   schema,
				Expected: mongoke.Map{"insertUser": mongoke.Map{"returning": mongoke.Map{"name": "xxx", "age": 10, "_id": "000000000000000000000000"}, "affectedCount": 1}},
				Query: `
				mutation {
					insertUser(data: {name: "xxx", age: 10, _id: "000000000000000000000000"}) {
						affectedCount
						returning {
							name
							age
							_id
						}
					}
				}
				`,
			},
			{
				Name:          "insert with missing required fields should error",
				Schema:        schema,
				Expected:      mongoke.Map{"insertUser": mongoke.Map{"returning": mongoke.Map{"name": "xxx", "age": 10, "_id": "000000000000000000000000"}, "affectedCount": 1}},
				ExpectedError: true,
				Query: `
				mutation {
					insertUser(data: {name: "xxx"}) {
						affectedCount
						returning {
							name
							age
							_id
						}
					}
				}
				`,
			},
		},
	})
}

func TestMutationWithMockedDb(t *testing.T) {
	typeDefs := `
	scalar ObjectId
	interface Named {
		name: String
	}

	type User implements Named {
		_id: ObjectId
		name: String
		age: Int
	}
	`
	typesConf := map[string]*mongoke.TypeConfig{
		"User": {Collection: "users"},
	}

	testutil.NewTestGroup(t, testutil.NewTestGroupParams{
		Collection: "users",
		Documents:  []mongoke.Map{},
		Tests: []testutil.TestCase{
			{
				Name: "insertone returns error",
				Schema: takeFirst(mongoke_schema.MakeMongokeSchema(mongoke.Config{
					Schema: typeDefs,
					DatabaseFunctions: &mock.DatabaseInterfaceMock{
						InsertManyFunc: func(ctx context.Context, p mongoke.InsertManyParams) (mongoke.NodesMutationPayload, error) {
							return mongoke.NodesMutationPayload{}, errors.New("error")
						},
					},
					Types: typesConf,
				})),
				// Expected:      mongoke.Map{"insertUser": mongoke.Map{"returning": nil, "affectedCount": 0}},
				ExpectedError: true,
				Query: `
				mutation {
					insertUser(data: {name: "xxx"}) {
						affectedCount
						returning {
							name
							age
							_id
						}
					}
				}
				`,
			},
		},
	})
}

func takeFirst(x, y interface{}) graphql.Schema {
	// t.Error(y)
	return x.(graphql.Schema)
}

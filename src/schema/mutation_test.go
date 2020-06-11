package schema_test

import (
	"context"
	"errors"
	"testing"

	"github.com/graphql-go/graphql"
	goke "github.com/remorses/goke/src"
	"github.com/remorses/goke/src/fakedata"
	"github.com/remorses/goke/src/mock"
	goke_schema "github.com/remorses/goke/src/schema"
	"github.com/remorses/goke/src/testutil"
)

func TestMutationWithEmptyFakeDatabase(t *testing.T) {
	db := &fakedata.FakeDatabaseFunctions{}
	schema, _ := goke_schema.MakeGokeSchema(goke.Config{
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
		Types: map[string]*goke.TypeConfig{
			"User": {Collection: "users"},
		},
	})

	testutil.NewTestGroup(t, testutil.NewTestGroupParams{
		Collection:    "users",
		Database:      db,
		Documents:     []goke.Map{},
		DefaultSchema: schema,
		Tests: []testutil.TestCase{
			{
				Name:     "update with set",
				Schema:   schema,
				Expected: goke.Map{"updateUser": goke.Map{"returning": nil, "affectedCount": 0}},
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
				Expected: goke.Map{"updateUser": goke.Map{"returning": nil, "affectedCount": 0}},
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
				Expected: goke.Map{"updateUserNodes": goke.Map{"returning": []goke.Map{}, "affectedCount": 0}},
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
				Expected: goke.Map{"insertUser": goke.Map{"returning": goke.Map{"name": "xxx", "age": 10, "_id": "000000000000000000000000"}, "affectedCount": 1}},
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
				Expected:      goke.Map{"insertUser": goke.Map{"returning": goke.Map{"name": "xxx", "age": 10, "_id": "000000000000000000000000"}, "affectedCount": 1}},
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
	typesConf := map[string]*goke.TypeConfig{
		"User": {Collection: "users"},
	}

	testutil.NewTestGroup(t, testutil.NewTestGroupParams{
		Collection: "users",
		Documents:  []goke.Map{},
		Tests: []testutil.TestCase{
			{
				Name: "insertone returns error",
				Schema: takeFirst(goke_schema.MakeGokeSchema(goke.Config{
					Schema: typeDefs,
					DatabaseFunctions: &mock.DatabaseInterfaceMock{
						InsertManyFunc: func(ctx context.Context, p goke.InsertManyParams) (goke.NodesMutationPayload, error) {
							return goke.NodesMutationPayload{}, errors.New("error")
						},
					},
					Types: typesConf,
				})),
				// Expected:      goke.Map{"insertUser": goke.Map{"returning": nil, "affectedCount": 0}},
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

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
				Name:     "updateOne_with_set",
				Schema:   schema,
				Expected: goke.Map{"updateOneUser": goke.Map{"returning": nil, "affectedCount": 0}},
				Query: `
				mutation {
					updateOneUser(set: {name: "xxx"}) {
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
				Name:     "updateOne with set and where",
				Schema:   schema,
				Expected: goke.Map{"updateOneUser": goke.Map{"returning": nil, "affectedCount": 0}},
				Query: `
				mutation {
					updateOneUser(where: {name: {eq: "zzz"}}, set: {name: "xxx"}) {
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
				Name:     "updateMany with set",
				Schema:   schema,
				Expected: goke.Map{"updateManyUser": goke.Map{"returning": []goke.Map{}, "affectedCount": 0}},
				Query: `
				mutation {
					updateManyUser(set: {name: "xxx"}) {
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
				Name:     "insertOne",
				Schema:   schema,
				Expected: goke.Map{"insertOneUser": goke.Map{"returning": goke.Map{"name": "xxx", "age": 10, "_id": "000000000000000000000000"}, "affectedCount": 1}},
				Query: `
				mutation {
					insertOneUser(data: {name: "xxx", age: 10, _id: "000000000000000000000000"}) {
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
				Name:          "insertOne with missing required fields should error",
				Schema:        schema,
				Expected:      goke.Map{"insertOneUser": goke.Map{"returning": goke.Map{"name": "xxx", "age": 10, "_id": "000000000000000000000000"}, "affectedCount": 1}},
				ExpectedError: true,
				Query: `
				mutation {
					insertOneUser(data: {name: "xxx"}) {
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

func TestDeleteWithEmptyFakeDatabase(t *testing.T) {
	exampleUsers := []goke.Map{
		{"name": "01", "age": 1},
		{"name": "02", "age": 2},
		{"name": "03", "age": 3},
	}
	db := &fakedata.FakeDatabaseFunctions{}
	schema, _ := goke_schema.MakeGokeSchema(goke.Config{
		Schema: `
		scalar ObjectId

		type User {
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
		Documents:     exampleUsers,
		DefaultSchema: schema,
		Tests: []testutil.TestCase{
			{
				Name:     "deleteMany removes everything",
				Schema:   schema,
				Expected: goke.Map{"deleteManyUser": goke.Map{"returning": exampleUsers, "affectedCount": len(exampleUsers)}},
				Query: `
				mutation {
					deleteManyUser {
						returning {
							name
							age
						}
						affectedCount
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
				Name: "insertOne returns error",
				Schema: takeFirst(goke_schema.MakeGokeSchema(goke.Config{
					Schema: typeDefs,
					DatabaseFunctions: &mock.DatabaseInterfaceMock{
						InsertManyFunc: func(ctx context.Context, p goke.InsertManyParams, h goke.TransformDocument) (goke.NodesMutationPayload, error) {
							return goke.NodesMutationPayload{}, errors.New("error")
						},
					},
					Types: typesConf,
				})),
				// Expected:      goke.Map{"insertOneUser": goke.Map{"returning": nil, "affectedCount": 0}},
				ExpectedError: true,
				Query: `
				mutation {
					insertOneUser(data: {name: "xxx"}) {
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

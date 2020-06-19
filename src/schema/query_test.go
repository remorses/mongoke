package schema_test

import (
	"fmt"
	"testing"

	"github.com/dgrijalva/jwt-go"
	goke "github.com/remorses/goke/src"
	"github.com/remorses/goke/src/fakedata"
	"github.com/remorses/goke/src/mongodb"
	goke_schema "github.com/remorses/goke/src/schema"
	"github.com/remorses/goke/src/testutil"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	False = false
	True  = true
)

func TestQueryReturnValuesWithMongoDB(t *testing.T) {
	exampleUsers := []goke.Map{
		{"name": "01", "age": 1},
		{"name": "02", "age": 2},
		{"name": "03", "age": 3},
	}
	for i, u := range exampleUsers {
		id, err := primitive.ObjectIDFromHex("00000000000000000000000" + fmt.Sprintf("%d", i))
		if err != nil {
			t.Error(err)
		}
		u["_id"] = id
	}
	mongo := &mongodb.MongodbDatabaseFunctions{
		Config: goke.Config{
			Mongodb: goke.MongodbConfig{
				Uri: testutil.MONGODB_URI,
			},
		},
	}
	exampleUser := exampleUsers[2]
	schema, _ := goke_schema.MakeGokeSchema(goke.Config{
		Schema: `
		scalar ObjectId
		interface Named {
			name: String
		}

		type User implements Named {
			_id: ObjectId
			name: String
			age: Int
		}
		`,
		DatabaseFunctions: mongo,
		Types: map[string]*goke.TypeConfig{
			"User": {Collection: "users"},
		},
	})
	testutil.NewTestGroup(t, testutil.NewTestGroupParams{
		Collection:    "users",
		Database:      mongo,
		Documents:     exampleUsers,
		DefaultSchema: schema,
		Tests: []testutil.TestCase{
			{
				Name:     "findOne query without args",
				Schema:   schema,
				Expected: goke.Map{"findOneUser": exampleUser},
				Query: `
				{
					findOneUser {
						name
						age
						_id
					}
				}
			`,
			},
			{
				Name:     "findOne query with eq",
				Schema:   schema,
				Expected: goke.Map{"findOneUser": goke.Map{"name": "03"}},
				Query: `
				{
					findOneUser(where: {name: {eq: "03"}}) {
						name
					}
				}
			`,
			},
			{
				Name:     "findOne query with eq objectId",
				Schema:   schema,
				Expected: goke.Map{"findOneUser": exampleUsers[0]},
				Query: `
				{
					findOneUser(where: {_id: {eq: "000000000000000000000000"}}) {
						_id
						age
						name
					}
				}
			`,
			},
			{
				Name:     "findMany simple",
				Schema:   schema,
				Expected: goke.Map{"findManyUser": exampleUsers},
				Query: `
				{
					findManyUser {
						_id
						name
						age
					}
				}
				`,
			},
			{
				Name:     "relay query without args",
				Schema:   schema,
				Expected: goke.Map{"UserNodes": goke.Map{"nodes": exampleUsers}},
				Query: `
				{
					UserNodes(direction: ASC) {
						nodes {
							_id
							name
							age
						}
					}
				}
				`,
			},

			{
				Name:     "relay query with first",
				Schema:   schema,
				Expected: goke.Map{"UserNodes": goke.Map{"nodes": exampleUsers[:2]}},
				Query: `
				{
					UserNodes(first: 2, direction: ASC) {
						nodes {
							_id
							name
							age
						}
					}
				}
				`,
			},
			{
				Name:     "relay query with last",
				Schema:   schema,
				Expected: goke.Map{"UserNodes": goke.Map{"nodes": exampleUsers[len(exampleUsers)-2:]}},
				Query: `
				{
					UserNodes(last: 2, direction: ASC) {
						nodes {
							_id
							name
							age
						}
					}
				}
			`,
			},
			{
				Name:   "relay query with string cursorField",
				Schema: schema,
				Expected: goke.Map{"UserNodes": goke.Map{
					"nodes":    exampleUsers[:2],
					"pageInfo": goke.Map{"endCursor": "02"},
				}},
				Query: `
				{
					UserNodes(first: 2, cursorField: name, direction: ASC) {
						nodes {
							_id
							name
							age
						}
						pageInfo {
							endCursor
						}
					}
				}
			`,
			},
			{
				Name:   "relay query with int cursorField",
				Schema: schema,
				Expected: goke.Map{"UserNodes": goke.Map{
					"nodes":    exampleUsers[:2],
					"pageInfo": goke.Map{"endCursor": 2},
				}},
				Query: `
				{
					UserNodes(first: 2, cursorField: age, direction: ASC) {
						nodes {
							_id
							name
							age
						}
						pageInfo {
							endCursor
						}
					}
				}
			`,
			},
			{
				Name:   "relay query with ObjectId cursorField",
				Schema: schema,
				Expected: goke.Map{"UserNodes": goke.Map{
					"nodes":    exampleUsers[:2],
					"pageInfo": goke.Map{"endCursor": exampleUsers[1]["_id"].(primitive.ObjectID).Hex()},
				}},
				Query: `
				{
					UserNodes(first: 2, direction: ASC) {
						nodes {
							_id
							name
							age
						}
						pageInfo {
							endCursor
						}
					}
				}
			`,
			},
		},
	})

}

func TestQueryWithFakeDatabase(t *testing.T) {
	exampleUsers := []goke.Map{
		{"name": "01", "age": 1},
		{"name": "02", "age": 2},
		{"name": "03", "age": 3},
	}
	for i, u := range exampleUsers {
		id, err := primitive.ObjectIDFromHex("00000000000000000000000" + fmt.Sprintf("%d", i))
		if err != nil {
			t.Error(err)
		}
		u["_id"] = id
	}
	db := &fakedata.FakeDatabaseFunctions{}
	exampleUser := exampleUsers[2]
	schema, _ := goke_schema.MakeGokeSchema(goke.Config{
		Schema: `
		scalar ObjectId
		interface Named {
			name: String
		}

		type User implements Named {
			_id: ObjectId
			name: String!
			age: Int
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
				Name:     "findOne query without args",
				Schema:   schema,
				Expected: goke.Map{"findOneUser": exampleUser},
				Query: `
			{
				findOneUser {
					name
					age
					_id
				}
			}
			`,
			},
			{
				Name:     "findOne query with eq",
				Schema:   schema,
				Expected: goke.Map{"findOneUser": goke.Map{"name": "03"}},
				Query: `
			{
				findOneUser(where: {name: {eq: "03"}}) {
					name
				}
			}
			`,
			},
			{
				Name:     "findOne query with eq objectId",
				Schema:   schema,
				Expected: goke.Map{"findOneUser": exampleUsers[0]},
				Query: `
			{
				findOneUser(where: {_id: {eq: "000000000000000000000000"}}) {
					_id
					age
					name
				}
			}
			`,
			},
			{
				Name:     "findMany query without args",
				Schema:   schema,
				Expected: goke.Map{"UserNodes": goke.Map{"nodes": exampleUsers}},
				Query: `
			{
				UserNodes(direction: ASC) {
					nodes {
						_id
						name
						age
					}
				}
			}
			`,
			},
			{
				Name:     "findMany simple",
				Schema:   schema,
				Expected: goke.Map{"findManyUser": exampleUsers},
				Query: `
			{
				findManyUser {
					_id
					name
					age
				}
			}
			`,
			},
			{
				Name:     "relay query with first",
				Schema:   schema,
				Expected: goke.Map{"UserNodes": goke.Map{"nodes": exampleUsers[:2]}},
				Query: `
			{
				UserNodes(first: 2, direction: ASC) {
					nodes {
						_id
						name
						age
					}
				}
			}
			`,
			},
			{
				Name:     "relay query with last",
				Schema:   schema,
				Expected: goke.Map{"UserNodes": goke.Map{"nodes": exampleUsers[len(exampleUsers)-2:]}},
				Query: `
			{
				UserNodes(last: 2, direction: ASC) {
					nodes {
						_id
						name
						age
					}
				}
			}
			`,
			},
			{
				Name:   "relay query with string cursorField",
				Schema: schema,
				Expected: goke.Map{"UserNodes": goke.Map{
					"nodes":    exampleUsers[:2],
					"pageInfo": goke.Map{"endCursor": "02"},
				}},
				Query: `
			{
				UserNodes(first: 2, cursorField: name, direction: ASC) {
					nodes {
						_id
						name
						age
					}
					pageInfo {
						endCursor
					}
				}
			}
			`,
			},
			{
				Name:   "relay query with int cursorField",
				Schema: schema,
				Expected: goke.Map{"UserNodes": goke.Map{
					"nodes":    exampleUsers[:2],
					"pageInfo": goke.Map{"endCursor": 2},
				}},
				Query: `
			{
				UserNodes(first: 2, cursorField: age, direction: ASC) {
					nodes {
						_id
						name
						age
					}
					pageInfo {
						endCursor
					}
				}
			}
			`,
			},
			{
				Name:   "relay query with ObjectId cursorField",
				Schema: schema,
				Expected: goke.Map{"UserNodes": goke.Map{
					"nodes":    exampleUsers[:2],
					"pageInfo": goke.Map{"endCursor": exampleUsers[1]["_id"].(primitive.ObjectID).Hex()},
				}},
				Query: `
			{
				UserNodes(first: 2, direction: ASC) {
					nodes {
						_id
						name
						age
					}
					pageInfo {
						endCursor
					}
				}
			}
			`,
			},
		},
	})

}

func TestToOneRelationWithFakeDatabase(t *testing.T) {
	exampleUsers := []goke.Map{
		{"name": "01", "age": 1},
		{"name": "02", "age": 2, "friendName": "01"},
		{"name": "03", "age": 3},
	}
	db := &fakedata.FakeDatabaseFunctions{}
	// exampleUser := exampleUsers[2]
	schema, _ := goke_schema.MakeGokeSchema(goke.Config{
		Schema: `
		scalar ObjectId

		type User {
			_id: ObjectId
			name: String!
			friendName: String
			age: Int
		}
		`,
		DatabaseFunctions: db,
		Types: map[string]*goke.TypeConfig{
			"User": {Collection: "users"},
		},
		Relations: []goke.RelationConfig{
			{
				Field:        "friend",
				From:         "User",
				To:           "User",
				RelationType: "to_one",
				Where: map[string]goke.Filter{
					"friendName": {
						Eq: "parent.name",
					},
				},
			},
		},
	})
	testutil.NewTestGroup(t, testutil.NewTestGroupParams{
		Collection:    "users",
		Database:      db,
		Documents:     exampleUsers,
		DefaultSchema: schema,
		Tests: []testutil.TestCase{

			{
				Name:   "findOne",
				Schema: schema,
				Expected: goke.Map{"findOneUser": goke.Map{
					"name": exampleUsers[0]["name"],
					"friend": goke.Map{
						"name": exampleUsers[1]["name"],
					},
				}},
				Query: `
			{
				findOneUser(where: {name: {eq: "01"}}) {
					name
					friend {
						name
					}
				}
			}
			`,
			},
		},
	})

}

func TestToManyRelationWithFakeDatabase(t *testing.T) {
	exampleUsers := []goke.Map{
		{"name": "01", "age": 1},
		{"name": "02", "age": 2, "friendName": "01"},
		{"name": "03", "age": 3, "friendName": "01"},
	}
	db := &fakedata.FakeDatabaseFunctions{}
	// exampleUser := exampleUsers[2]
	schema, _ := goke_schema.MakeGokeSchema(goke.Config{
		Schema: `
		scalar ObjectId

		type User {
			_id: ObjectId
			name: String!
			friendName: String
			age: Int
		}
		`,
		DatabaseFunctions: db,
		Types: map[string]*goke.TypeConfig{
			"User": {Collection: "users"},
		},
		Relations: []goke.RelationConfig{
			{
				Field:        "friends",
				From:         "User",
				To:           "User",
				RelationType: "to_many",
				Where: map[string]goke.Filter{
					"friendName": {
						Eq: "parent.name",
					},
				},
			},
		},
	})
	testutil.NewTestGroup(t, testutil.NewTestGroupParams{
		Collection:    "users",
		Database:      db,
		Documents:     exampleUsers,
		DefaultSchema: schema,
		Tests: []testutil.TestCase{

			{
				Name:   "findOne relation",
				Schema: schema,
				Expected: goke.Map{"findOneUser": goke.Map{
					"name": exampleUsers[0]["name"],
					"friends": []goke.Map{
						{"name": exampleUsers[1]["name"]},
						{"name": exampleUsers[2]["name"]},
					},
				}},
				Query: `
			{
				findOneUser(where: {name: {eq: "01"}}) {
					name
					friends {
						name
					}
				}
			}
			`,
			},
		},
	})

}

func TestQueryWithEmptyFakeDatabase(t *testing.T) {
	db := &fakedata.FakeDatabaseFunctions{}
	emptyList := []goke.Map{}
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
		Documents:     emptyList,
		DefaultSchema: schema,
		Tests: []testutil.TestCase{
			{
				Name:     "findOne query without args",
				Schema:   schema,
				Expected: goke.Map{"findOneUser": nil},
				Query: `
			{
				findOneUser {
					name
					age
					_id
				}
			}
			`,
			},
			{
				Name:     "findOne query with eq",
				Schema:   schema,
				Expected: goke.Map{"findOneUser": nil},
				Query: `
			{
				findOneUser(where: {name: {eq: "03"}}) {
					name
				}
			}
			`,
			},
			{
				Name:     "findMany with list",
				Schema:   schema,
				Expected: goke.Map{"findManyUser": emptyList},
				Query: `
			{
				findManyUser {
					_id
					name
					age
				}
			}
			`,
			},
			{
				Name:     "findMany query with first",
				Schema:   schema,
				Expected: goke.Map{"UserNodes": goke.Map{"nodes": emptyList}},
				Query: `
			{
				UserNodes(first: 2, direction: ASC) {
					nodes {
						_id
						name
						age
					}
				}
			}
			`,
			},

			{
				Name:   "findMany query with string cursorField",
				Schema: schema,
				Expected: goke.Map{"UserNodes": goke.Map{
					"nodes":    emptyList,
					"pageInfo": goke.Map{"endCursor": nil},
				}},
				Query: `
			{
				UserNodes(first: 2, cursorField: name, direction: ASC) {
					nodes {
						_id
						name
						age
					}
					pageInfo {
						endCursor
					}
				}
			}
			`,
			},
		},
	})

}

func TestAdminCheck(t *testing.T) {
	exampleUsers := []goke.Map{
		{"name": "01", "age": 1},
		{"name": "02", "age": 2},
		{"name": "03", "age": 3},
	}
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
		DefaultPermissions: []string{},
	})

	testutil.NewTestGroup(t, testutil.NewTestGroupParams{
		Collection:    "users",
		Database:      db,
		Documents:     exampleUsers,
		DefaultSchema: schema,
		Tests: []testutil.TestCase{
			{
				Name:   "findMany without admin is blocked",
				Schema: schema,
				RootObject: goke.Map{
					"isAdmin": false,
				},
				ExpectedError: true, // cannot execute read operation with current user permissions
				Query: `
				{
					findManyUser {
						name
						age
						_id
					}
				}
			`,
			},
			{
				Name:   "findMany with admin skip checks",
				Schema: schema,
				RootObject: goke.Map{
					"isAdmin": true,
				},
				ExpectedError: true,
				Expected:      goke.Map{"findManyUser": exampleUsers},
				Query: `
				{
					findManyUser {
						name
						age
						_id
					}
				}
			`,
			},
		},
	})

}

func TestReadPermissions(t *testing.T) {
	exampleUsers := []goke.Map{
		{"name": "01", "age": 1},
		{"name": "02", "age": 2},
		{"name": "03", "age": 3},
	}
	db := &fakedata.FakeDatabaseFunctions{}
	typeDefs := `
	scalar ObjectId
	interface Named {
		name: String
	}

	type User implements Named {
		_id: ObjectId!
		name: String
		age: Int!
	}
	`
	userId := "0000000"
	s, err := goke_schema.MakeGokeSchema(goke.Config{
		Schema:            typeDefs,
		DatabaseFunctions: db,
		Types: map[string]*goke.TypeConfig{
			"User": {
				Collection: "users",
				Permissions: []goke.AuthGuard{
					{
						Expression:        fmt.Sprintf("jwt.user_id == \"%s\"", userId),
						AllowedOperations: []string{goke.Operations.READ},
					},
				},
			},
		},
		DefaultPermissions: []string{},
	})
	if err != nil {
		t.Error(err)
	}

	testutil.NewTestGroup(t, testutil.NewTestGroupParams{
		Collection: "users",
		Database:   db,
		Documents:  exampleUsers,
		Tests: []testutil.TestCase{
			{
				Name:          "findMany without jwt errors",
				Schema:        s,
				ExpectedError: true, // cannot execute read operation with current user permissions
				Query: `
				{
					findManyUser {
						name
						age
						_id
					}
				}
			`,
			},
			{
				Name:     "findMany with jwt expressin passes",
				Schema:   s,
				Expected: goke.Map{"findManyUser": exampleUsers},
				RootObject: goke.Map{
					"jwt": jwt.MapClaims{
						"user_id": userId,
					},
				},
				Query: `
				{
					findManyUser {
						name
						age
					}
				}
			`,
			},
		},
	})

}

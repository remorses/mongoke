package schema_test

import (
	"fmt"
	"testing"

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

package schema

import (
	"fmt"
	"testing"

	mongoke "github.com/remorses/mongoke/src"
	"github.com/remorses/mongoke/src/fakedata"
	"github.com/remorses/mongoke/src/mongodb"
	"github.com/remorses/mongoke/src/testutil"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	False = false
	True  = true
)

func TestQueryReturnValuesWithMongoDB(t *testing.T) {
	exampleUsers := []mongoke.Map{
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
		Config: mongoke.Config{
			Mongodb: mongoke.MongodbConfig{
				Uri: testutil.MONGODB_URI,
			},
		},
	}
	exampleUser := exampleUsers[2]
	schema, _ := MakeMongokeSchema(mongoke.Config{
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
		Types: map[string]*mongoke.TypeConfig{
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
				Expected: mongoke.Map{"User": exampleUser},
				Query: `
			{
				User {
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
				Expected: mongoke.Map{"User": mongoke.Map{"name": "03"}},
				Query: `
			{
				User(where: {name: {eq: "03"}}) {
					name
				}
			}
			`,
			},
			{
				Name:     "findOne query with eq objectId",
				Schema:   schema,
				Expected: mongoke.Map{"User": exampleUsers[0]},
				Query: `
			{
				User(where: {_id: {eq: "000000000000000000000000"}}) {
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
				Expected: mongoke.Map{"UserNodes": mongoke.Map{"nodes": exampleUsers}},
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
				Name:     "findMany with list",
				Schema:   schema,
				Expected: mongoke.Map{"UserList": exampleUsers},
				Query: `
			{
				UserList {
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
				Expected: mongoke.Map{"UserNodes": mongoke.Map{"nodes": exampleUsers[:2]}},
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
				Name:     "findMany query with last",
				Schema:   schema,
				Expected: mongoke.Map{"UserNodes": mongoke.Map{"nodes": exampleUsers[len(exampleUsers)-2:]}},
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
				Name:   "findMany query with string cursorField",
				Schema: schema,
				Expected: mongoke.Map{"UserNodes": mongoke.Map{
					"nodes":    exampleUsers[:2],
					"pageInfo": mongoke.Map{"endCursor": "02"},
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
				Name:   "findMany query with int cursorField",
				Schema: schema,
				Expected: mongoke.Map{"UserNodes": mongoke.Map{
					"nodes":    exampleUsers[:2],
					"pageInfo": mongoke.Map{"endCursor": 2},
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
				Name:   "findMany query with ObjectId cursorField",
				Schema: schema,
				Expected: mongoke.Map{"UserNodes": mongoke.Map{
					"nodes":    exampleUsers[:2],
					"pageInfo": mongoke.Map{"endCursor": exampleUsers[1]["_id"].(primitive.ObjectID).Hex()},
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
	exampleUsers := []mongoke.Map{
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
	schema, _ := MakeMongokeSchema(mongoke.Config{
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
		Types: map[string]*mongoke.TypeConfig{
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
				Expected: mongoke.Map{"User": exampleUser},
				Query: `
			{
				User {
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
				Expected: mongoke.Map{"User": mongoke.Map{"name": "03"}},
				Query: `
			{
				User(where: {name: {eq: "03"}}) {
					name
				}
			}
			`,
			},
			{
				Name:     "findOne query with eq objectId",
				Schema:   schema,
				Expected: mongoke.Map{"User": exampleUsers[0]},
				Query: `
			{
				User(where: {_id: {eq: "000000000000000000000000"}}) {
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
				Expected: mongoke.Map{"UserNodes": mongoke.Map{"nodes": exampleUsers}},
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
				Name:     "findMany with list",
				Schema:   schema,
				Expected: mongoke.Map{"UserList": exampleUsers},
				Query: `
			{
				UserList {
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
				Expected: mongoke.Map{"UserNodes": mongoke.Map{"nodes": exampleUsers[:2]}},
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
				Name:     "findMany query with last",
				Schema:   schema,
				Expected: mongoke.Map{"UserNodes": mongoke.Map{"nodes": exampleUsers[len(exampleUsers)-2:]}},
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
				Name:   "findMany query with string cursorField",
				Schema: schema,
				Expected: mongoke.Map{"UserNodes": mongoke.Map{
					"nodes":    exampleUsers[:2],
					"pageInfo": mongoke.Map{"endCursor": "02"},
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
				Name:   "findMany query with int cursorField",
				Schema: schema,
				Expected: mongoke.Map{"UserNodes": mongoke.Map{
					"nodes":    exampleUsers[:2],
					"pageInfo": mongoke.Map{"endCursor": 2},
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
				Name:   "findMany query with ObjectId cursorField",
				Schema: schema,
				Expected: mongoke.Map{"UserNodes": mongoke.Map{
					"nodes":    exampleUsers[:2],
					"pageInfo": mongoke.Map{"endCursor": exampleUsers[1]["_id"].(primitive.ObjectID).Hex()},
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
	emptyList := []mongoke.Map{}
	schema, _ := MakeMongokeSchema(mongoke.Config{
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
		Documents:     emptyList,
		DefaultSchema: schema,
		Tests: []testutil.TestCase{
			{
				Name:     "findOne query without args",
				Schema:   schema,
				Expected: mongoke.Map{"User": nil},
				Query: `
			{
				User {
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
				Expected: mongoke.Map{"User": nil},
				Query: `
			{
				User(where: {name: {eq: "03"}}) {
					name
				}
			}
			`,
			},
			{
				Name:     "findMany with list",
				Schema:   schema,
				Expected: mongoke.Map{"UserList": emptyList},
				Query: `
			{
				UserList {
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
				Expected: mongoke.Map{"UserNodes": mongoke.Map{"nodes": emptyList}},
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
				Expected: mongoke.Map{"UserNodes": mongoke.Map{
					"nodes":    emptyList,
					"pageInfo": mongoke.Map{"endCursor": nil},
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

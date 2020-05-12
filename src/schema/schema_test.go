package schema

import (
	"context"
	"testing"

	"github.com/go-test/deep"
	mongoke "github.com/remorses/mongoke/src"
	"github.com/remorses/mongoke/src/mock"
	"github.com/remorses/mongoke/src/mongodb"
	"github.com/remorses/mongoke/src/testutil"
	"github.com/stretchr/testify/require"
)

var False = false
var True = true

func TestQueryReturnValues(t *testing.T) {
	var exampleUsers = []mongoke.Map{
		{"name": "01", "age": 1},
		{"name": "02", "age": 2},
		{"name": "03", "age": 3},
	}
	exampleUser := exampleUsers[0]
	databaseMock := &mock.DatabaseInterfaceMock{
		FindManyFunc: func(p mongoke.FindManyParams) ([]mongoke.Map, error) {
			return exampleUsers, nil
		},
		FindOneFunc: func(p mongoke.FindOneParams) (interface{}, error) {
			return exampleUser, nil
		},
	}
	var config = mongoke.Config{
		Schema: `
		interface Named {
			name: String
		}

		type User implements Named {
			name: String
			age: Int
		}
		`,
		DatabaseUri: testutil.MONGODB_URI,
		Types: map[string]*mongoke.TypeConfig{
			"User": {Collection: "users"},
		},
		DatabaseFunctions: databaseMock,
	}

	cases := []struct {
		Name          string
		Query         string
		Expected      mongoke.Map
		ExpectedError string
		Config        mongoke.Config
	}{
		{
			Name:     "findOne query without args",
			Config:   config,
			Expected: mongoke.Map{"User": exampleUser},
			Query: `
			{
				User {
					name
					age
				}
			}
			`,
		},
		{
			Name: "findOne query with permissions false",
			Config: mongoke.Config{
				Schema: `
				type User {
					name: String
					age: Int
				}
				`,
				DatabaseUri: testutil.MONGODB_URI,
				Types: map[string]*mongoke.TypeConfig{
					"User": {
						Collection: "users",
						Permissions: []mongoke.AuthGuard{
							{
								Expression: "false",
							},
						},
					},
				},
				DatabaseFunctions: databaseMock,
			},
			ExpectedError: "no enough permissions", // TODO tests should check errors name
			Query: `
			{
				User {
					name
					age
				}
			}
			`,
		},
		{
			Name: "schema with Union type",
			Config: mongoke.Config{
				Schema: `
				type Guest {
					name: String
				}
				type Admin {
					password: String
				}
				union User = Admin | Guest
				`,
				DatabaseUri: testutil.MONGODB_URI,
				Types: map[string]*mongoke.TypeConfig{
					"Admin": {
						Exposed:  &False,
						IsTypeOf: "false",
					},
					"Guest": {
						Exposed:  &False,
						IsTypeOf: "x.name == \"01\"",
					},
					"User": {
						Collection: "users",
					},
				},
				DatabaseFunctions: databaseMock,
			},
			Expected: mongoke.Map{"User": mongoke.Map{"name": "01"}},
			Query: `
			{
				User {
					...on Guest {
						name
					}
					...on Admin {
						password
					}
				}
			}
			`,
		},
		{
			Name: "findOne query with permissions HideFields",
			Config: mongoke.Config{
				Schema: `
				type User {
					name: String
					age: Int
				}
				`,
				DatabaseUri: testutil.MONGODB_URI,
				Types: map[string]*mongoke.TypeConfig{
					"User": {
						Collection: "users",

						Permissions: []mongoke.AuthGuard{
							{
								Expression: "true",
								HideFields: []string{"name"},
							},
						},
					},
				},
				DatabaseFunctions: databaseMock,
			},
			Expected: mongoke.Map{"User": mongoke.Map{"name": nil, "age": 1}},
			Query: `
			{
				User {
					name
					age
				}
			}
			`,
		},
		{
			Name:     "findMany query without args",
			Config:   config,
			Expected: mongoke.Map{"UserNodes": mongoke.Map{"nodes": exampleUsers}},
			Query: `
			{
				UserNodes {
					nodes {
						name
						age
					}
				}
			}
			`,
		},
		{
			Name: "findMany query with permissions false",
			Config: mongoke.Config{
				Schema: `
				type User {
					name: String
					age: Int
				}
				`,
				DatabaseUri:       testutil.MONGODB_URI,
				DatabaseFunctions: databaseMock,
				Types: map[string]*mongoke.TypeConfig{
					"User": {
						Collection: "users",
						Permissions: []mongoke.AuthGuard{
							{
								Expression: "false",
							},
						},
					},
				},
			},
			Expected: mongoke.Map{"UserNodes": mongoke.Map{"nodes": []mongoke.Map{}}},
			Query: `
			{
				UserNodes {
					nodes {
						name
						age
					}
				}
			}
			`,
		},
		{
			Name:     "findOne query without args",
			Config:   config,
			Expected: mongoke.Map{"User": exampleUser},
			Query: `
			{
				User {
					name
					age
				}
			}
			`,
		},
		{
			Name: "findOne with to_many relation",
			Config: mongoke.Config{
				Schema: `
				type User {
					name: String
					age: Int
				}
				`,
				DatabaseUri:       testutil.MONGODB_URI,
				DatabaseFunctions: databaseMock,
				Types: map[string]*mongoke.TypeConfig{
					"User": {
						Collection: "users",
					},
				},
				Relations: []mongoke.RelationConfig{
					{
						Field:        "friends",
						From:         "User",
						To:           "User",
						RelationType: "to_many",
						Where:        make(map[string]mongoke.Filter),
					},
				},
			},
			Expected: mongoke.Map{"User": mongoke.Map{"name": "01", "friends": mongoke.Map{"nodes": exampleUsers}}},
			Query: `
			{
				User {
					name
					friends {
						nodes {
							name
							age
						}
					}
				}
			}
			`,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Log()
			// t.Log(testCase.Name)
			schema, err := MakeMongokeSchema(testCase.Config)
			if err != nil {
				t.Fatal(err)
			}
			if testCase.ExpectedError != "" {
				err = testutil.QuerySchemaShouldFail(t, schema, testCase.Query)
				return
			}
			res := testutil.QuerySchema(t, schema, testCase.Query)
			res = testutil.ConvertToPlainMap(res)
			expected := testutil.ConvertToPlainMap(testCase.Expected)
			t.Log("expected:", expected)
			t.Log("result:", res)
			require.Equal(t, testutil.Pretty(expected), testutil.Pretty(res))
		})
	}
}

func TestQueryReturnValuesWithMongoDB(t *testing.T) {
	collection := "users"
	var exampleUsers = []mongoke.Map{
		{"name": "01", "age": 1},
		{"name": "02", "age": 2},
		{"name": "03", "age": 3},
	}
	exampleUser := exampleUsers[0]
	var config = mongoke.Config{
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
		DatabaseUri: testutil.MONGODB_URI,
		Types: map[string]*mongoke.TypeConfig{
			"User": {Collection: "users"},
		},
	}

	cases := []struct {
		Name          string
		Query         string
		Expected      mongoke.Map
		ExpectedError string
		Config        mongoke.Config
	}{
		{
			Name:     "findOne query without args",
			Config:   config,
			Expected: mongoke.Map{"User": exampleUser},
			Query: `
			{
				User {
					name
					age
				}
			}
			`,
		},
		{
			Name:     "findOne query with eq",
			Config:   config,
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
			Name:     "findMany query without args",
			Config:   config,
			Expected: mongoke.Map{"UserNodes": mongoke.Map{"nodes": exampleUsers}},
			Query: `
			{
				UserNodes {
					nodes {
						name
						age
					}
				}
			}
			`,
		},
		{
			Name:     "findMany query with first",
			Config:   config,
			Expected: mongoke.Map{"UserNodes": mongoke.Map{"nodes": exampleUsers[:2]}},
			Query: `
			{
				UserNodes(first: 2) {
					nodes {
						name
						age
					}
				}
			}
			`,
		},
		{
			Name:     "findMany query with last",
			Config:   config,
			Expected: mongoke.Map{"UserNodes": mongoke.Map{"nodes": exampleUsers[len(exampleUsers)-2:]}},
			Query: `
			{
				UserNodes(last: 2) {
					nodes {
						name
						age
					}
				}
			}
			`,
		},
		{
			Name:     "findMany query with string cursorField",
			Config:   config,
			Expected: mongoke.Map{"UserNodes": mongoke.Map{"nodes": exampleUsers[:2], "pageInfo": mongoke.Map{"endCursor": "02"}}},
			Query: `
			{
				UserNodes(first: 2, cursorField: name) {
					nodes {
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
			Name:     "findMany query with int cursorField",
			Config:   config,
			Expected: mongoke.Map{"UserNodes": mongoke.Map{"nodes": exampleUsers[:2], "pageInfo": mongoke.Map{"endCursor": 2}}},
			Query: `
			{
				UserNodes(first: 2, cursorField: age) {
					nodes {
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
	}

	for _, testCase := range cases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Log()
			// t.Log(testCase.Name)
			schema, err := MakeMongokeSchema(testCase.Config)
			m := mongodb.MongodbDatabaseFunctions{}
			db, err := m.InitMongo(testutil.MONGODB_URI)
			if err != nil {
				t.Error(err)
			}
			// clear
			_, err = db.Collection(collection).DeleteMany(context.TODO(), mongoke.Map{})
			if err != nil {
				t.Error(err)
			}
			if err != nil {
				t.Error(err)
			}
			for _, user := range exampleUsers {
				_, err := db.Collection(collection).InsertOne(context.TODO(), user)
				if err != nil {
					t.Error(err)
				}
			}
			if testCase.ExpectedError != "" {
				err = testutil.QuerySchemaShouldFail(t, schema, testCase.Query)
				return
			}
			res := testutil.QuerySchema(t, schema, testCase.Query)
			res = testutil.ConvertToPlainMap(res)
			expected := testutil.ConvertToPlainMap(testCase.Expected)
			t.Log("expected:", expected)
			t.Log("result:", res)
			t.Log("expected:", testutil.Pretty(expected))
			t.Log("result:", testutil.Pretty(res))
			// require.Equal(t, testutil.Pretty(res), testutil.Pretty(expected))
			if diff := deep.Equal(res, expected); diff != nil {
				t.Error(diff)
			}
		})
	}
}

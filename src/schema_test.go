package mongoke

import (
	"context"
	"testing"

	"github.com/go-test/deep"
	"github.com/remorses/mongoke/src/testutil"
	"github.com/stretchr/testify/require"
)

func TestQueryArgs(t *testing.T) {
	databaseMock := &DatabaseInterfaceMock{
		FindManyFunc: func(p FindManyParams) ([]Map, error) {
			return nil, nil
		},
		FindOneFunc: func(p FindOneParams) (interface{}, error) {
			return nil, nil
		},
	}
	var config = Config{
		Schema: `
		type User {
			name: String
			age: Int
		}
		`,
		DatabaseUri: testutil.MONGODB_URI,
		Types: map[string]*TypeConfig{
			"User": {Collection: "users"},
		},
		databaseFunctions: databaseMock,
	}

	schema, err := MakeMongokeSchema(config)
	if err != nil {
		t.Error(err)
	}
	t.Run("introspect schema", func(t *testing.T) {
		if err != nil {
			t.Error(err)
		}
		testutil.QuerySchema(t, schema, testutil.IntrospectionQuery)
	})
	t.Run("findOne query without args", func(t *testing.T) {
		query := `
		{
			User {
				name
				age
			}
		}
		`
		testutil.QuerySchema(t, schema, query)
		calls := len(databaseMock.FindOneCalls())
		require.Equal(t, 1, calls)
		where := databaseMock.FindOneCalls()[calls-1].P.Where
		// require.Equal(t, nil, where)
		t.Log(where)
	})
	t.Run("findOne query with eq", func(t *testing.T) {
		databaseMock.calls.FindOne = nil
		query := `
		{
			User(where: {name: {eq: "xxx"}}) {
				name
				age
			}
		}
		`
		testutil.QuerySchema(t, schema, query)
		calls := len(databaseMock.FindOneCalls())
		require.Equal(t, 1, calls)
		where := databaseMock.FindOneCalls()[0].P.Where
		t.Log(pretty(where))
		require.Equal(t, "xxx", where["name"].Eq)
	})
	t.Run("findMany query with first, after, where", func(t *testing.T) {
		databaseMock.calls.FindMany = nil
		query := `
		{
			UserNodes(first: 10, after: "xxx", where: {name: {eq: "xxx"}}) {
				nodes {
					name
				}
			}
		}
		`
		testutil.QuerySchema(t, schema, query)
		calls := len(databaseMock.calls.FindMany)
		require.Equal(t, 1, calls)
		p := databaseMock.calls.FindMany[0].P
		t.Log("params", pretty(p))
		// + 1 because we need to know hasNextPage
		require.Equal(t, 10+1, p.Pagination.First)
		require.Equal(t, "xxx", p.Pagination.After)
	})
}

var trueValue = true
var falseValue = false

func TestQueryReturnValues(t *testing.T) {
	var exampleUsers = []Map{
		{"name": "01", "age": 1},
		{"name": "02", "age": 2},
		{"name": "03", "age": 3},
	}
	exampleUser := exampleUsers[0]
	databaseMock := &DatabaseInterfaceMock{
		FindManyFunc: func(p FindManyParams) ([]Map, error) {
			return exampleUsers, nil
		},
		FindOneFunc: func(p FindOneParams) (interface{}, error) {
			return exampleUser, nil
		},
	}
	var config = Config{
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
		Types: map[string]*TypeConfig{
			"User": {Collection: "users"},
		},
		databaseFunctions: databaseMock,
	}

	cases := []struct {
		Name          string
		Query         string
		Expected      Map
		ExpectedError string
		Config        Config
	}{
		{
			Name:     "findOne query without args",
			Config:   config,
			Expected: Map{"User": exampleUser},
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
			Config: Config{
				Schema: `
				type User {
					name: String
					age: Int
				}
				`,
				DatabaseUri: testutil.MONGODB_URI,
				Types: map[string]*TypeConfig{
					"User": {
						Collection: "users",
						Permissions: []AuthGuard{
							{
								Expression: "false",
							},
						},
					},
				},
				databaseFunctions: databaseMock,
			},
			ExpectedError: "no enough permissions", // TODO error not right
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
			Config: Config{
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
				Types: map[string]*TypeConfig{
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
				databaseFunctions: databaseMock,
			},
			Expected: Map{"User": Map{"name": "01"}},
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
			Config: Config{
				Schema: `
				type User {
					name: String
					age: Int
				}
				`,
				DatabaseUri: testutil.MONGODB_URI,
				Types: map[string]*TypeConfig{
					"User": {
						Collection: "users",

						Permissions: []AuthGuard{
							{
								Expression: "true",
								HideFields: []string{"name"},
							},
						},
					},
				},
				databaseFunctions: databaseMock,
			},
			Expected: Map{"User": Map{"name": nil, "age": 1}},
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
			Expected: Map{"UserNodes": Map{"nodes": exampleUsers}},
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
			Config: Config{
				Schema: `
				type User {
					name: String
					age: Int
				}
				`,
				DatabaseUri:       testutil.MONGODB_URI,
				databaseFunctions: databaseMock,
				Types: map[string]*TypeConfig{
					"User": {
						Collection: "users",
						Permissions: []AuthGuard{
							{
								Expression: "false",
							},
						},
					},
				},
			},
			Expected: Map{"UserNodes": Map{"nodes": []Map{}}},
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
			Expected: Map{"User": exampleUser},
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
			Config: Config{
				Schema: `
				type User {
					name: String
					age: Int
				}
				`,
				DatabaseUri:       testutil.MONGODB_URI,
				databaseFunctions: databaseMock,
				Types: map[string]*TypeConfig{
					"User": {
						Collection: "users",
					},
				},
				Relations: []RelationConfig{
					{
						Field:        "friends",
						From:         "User",
						To:           "User",
						RelationType: "to_many",
						Where:        make(map[string]Filter),
					},
				},
			},
			Expected: Map{"User": Map{"name": "01", "friends": Map{"nodes": exampleUsers}}},
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
			require.Equal(t, pretty(expected), pretty(res))
		})
	}
}

func TestQueryReturnValuesWithMongoDB(t *testing.T) {
	collection := "users"
	var exampleUsers = []Map{
		{"name": "01", "age": 1},
		{"name": "02", "age": 2},
		{"name": "03", "age": 3},
	}
	exampleUser := exampleUsers[0]
	var config = Config{
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
		Types: map[string]*TypeConfig{
			"User": {Collection: "users"},
		},
	}

	cases := []struct {
		Name          string
		Query         string
		Expected      Map
		ExpectedError string
		Config        Config
	}{
		{
			Name:     "findOne query without args",
			Config:   config,
			Expected: Map{"User": exampleUser},
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
			Expected: Map{"User": Map{"name": "03"}},
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
			Expected: Map{"UserNodes": Map{"nodes": exampleUsers}},
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
			Expected: Map{"UserNodes": Map{"nodes": exampleUsers[:2]}},
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
			Expected: Map{"UserNodes": Map{"nodes": exampleUsers[len(exampleUsers)-2:]}},
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
			Expected: Map{"UserNodes": Map{"nodes": exampleUsers[:2], "pageInfo": Map{"endCursor": "02"}}},
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
			Expected: Map{"UserNodes": Map{"nodes": exampleUsers[:2], "pageInfo": Map{"endCursor": 2}}},
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
			m := MongodbDatabaseFunctions{}
			db, err := m.initMongo(testutil.MONGODB_URI)
			if err != nil {
				t.Error(err)
			}
			// clear
			_, err = db.Collection(collection).DeleteMany(context.TODO(), Map{})
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
			t.Log("expected:", pretty(expected))
			t.Log("result:", pretty(res))
			// require.Equal(t, pretty(res), pretty(expected))
			if diff := deep.Equal(res, expected); diff != nil {
				t.Error(diff)
			}
		})
	}
}

package schema

import (
	"context"
	"testing"

	"github.com/go-test/deep"
	mongoke "github.com/remorses/mongoke/src"
	"github.com/remorses/mongoke/src/mongodb"
	"github.com/remorses/mongoke/src/testutil"
)

var (
	False = false
	True  = true
)

func TestQueryReturnValuesWithMongoDB(t *testing.T) {
	collection := "users"
	exampleUsers := []mongoke.Map{
		{"name": "01", "age": 1},
		{"name": "02", "age": 2},
		{"name": "03", "age": 3},
	}
	exampleUser := exampleUsers[2]
	config := mongoke.Config{
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
		DatabaseFunctions: mongodb.MongodbDatabaseFunctions{},
		DatabaseUri:       testutil.MONGODB_URI,
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
				UserNodes(direction: ASC) {
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
				UserNodes(first: 2, direction: ASC) {
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
				UserNodes(last: 2, direction: ASC) {
					nodes {
						name
						age
					}
				}
			}
			`,
		},
		{
			Name:   "findMany query with string cursorField",
			Config: config,
			Expected: mongoke.Map{"UserNodes": mongoke.Map{
				"nodes":    exampleUsers[:2],
				"pageInfo": mongoke.Map{"endCursor": "02"},
			}},
			Query: `
			{
				UserNodes(first: 2, cursorField: name, direction: ASC) {
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
			Name:   "findMany query with int cursorField",
			Config: config,
			Expected: mongoke.Map{"UserNodes": mongoke.Map{
				"nodes":    exampleUsers[:2],
				"pageInfo": mongoke.Map{"endCursor": 2},
			}},
			Query: `
			{
				UserNodes(first: 2, cursorField: age, direction: ASC) {
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
			ctx := context.Background()
			// t.Log(testCase.Name)
			schema, err := MakeMongokeSchema(testCase.Config)
			if err != nil {
				t.Error(err)
			}
			m := mongodb.MongodbDatabaseFunctions{}
			db, err := m.Init(ctx, testutil.MONGODB_URI)
			if err != nil {
				t.Error(err)
			}
			// clear
			_, err = db.Collection(collection).DeleteMany(ctx, mongoke.Map{})
			if err != nil {
				t.Error(err)
			}
			if err != nil {
				t.Error(err)
			}
			for _, user := range exampleUsers {
				_, err := db.Collection(collection).InsertOne(ctx, user)
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
			_, err = db.Collection(collection).DeleteMany(ctx, mongoke.Map{})
		})
	}

}

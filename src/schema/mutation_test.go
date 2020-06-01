package schema

import (
	"context"
	"testing"

	"github.com/go-test/deep"
	mongoke "github.com/remorses/mongoke/src"
	"github.com/remorses/mongoke/src/fakedata"
	"github.com/remorses/mongoke/src/testutil"
)

func TestMutationWithEmptyFakeDatabase(t *testing.T) {
	collection := "users"
	emptyList := []mongoke.Map{}
	fake := fakedata.FakeDatabaseFunctions{}
	config := mongoke.Config{
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
		DatabaseFunctions: &fake,
		Types: map[string]*mongoke.TypeConfig{
			"User": {Collection: "users"},
		},
	}

	cases := []struct {
		Name          string
		Query         string
		Expected      mongoke.Map
		ExpectedError bool
		Config        mongoke.Config
	}{
		{
			Name:     "updateOne with set",
			Config:   config,
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
			Name:     "updateOne with set and where",
			Config:   config,
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
			Name:     "insert",
			Config:   config,
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
			Config:        config,
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
			db, err := fake.Init(ctx)
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
			for _, user := range emptyList {
				_, err := db.Collection(collection).InsertOne(ctx, user)
				if err != nil {
					t.Error(err)
				}
			}
			if testCase.ExpectedError {
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

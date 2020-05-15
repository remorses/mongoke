package mongodb

import (
	"context"
	"testing"

	mongoke "github.com/remorses/mongoke/src"
	"github.com/remorses/mongoke/src/testutil"
	"github.com/stretchr/testify/require"
)

func TestFindMany(t *testing.T) {
	collection := "users"
	ctx := context.Background()
	uri := testutil.MONGODB_URI
	type Case struct {
		params   mongoke.FindManyParams
		expected []mongoke.Map
	}
	exampleUsers := []mongoke.Map{
		{"name": "01", "age": 1},
		{"name": "02", "age": 2},
		{"name": "03", "age": 3},
	}
	cases := map[string]Case{
		"no args": {
			params: mongoke.FindManyParams{
				Collection: collection,
			},
			expected: exampleUsers,
		},
		"Limit": {
			params: mongoke.FindManyParams{
				Collection: collection,
				Limit:      2,
			},
			expected: exampleUsers[:2],
		},
	}

	// clear and insert some docs
	m := MongodbDatabaseFunctions{}
	db, err := m.InitMongo(uri)
	if err != nil {
		t.Error(err)
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err = db.Collection(collection).DeleteMany(ctx, mongoke.Map{})
			if err != nil {
				t.Error(err)
			}
			for _, user := range exampleUsers {
				_, err := db.Collection(collection).InsertOne(ctx, user)
				if err != nil {
					t.Error(err)
				}
			}
			result, err := m.FindMany(c.params)
			if err != nil {
				t.Error(err)
			}
			t.Log("expected", testutil.Pretty(result))
			t.Log("result", testutil.Pretty(result))
			require.Equal(t, testutil.Pretty(c.expected), testutil.Pretty(result))
		})
	}
	_, err = db.Collection(collection).DeleteMany(ctx, mongoke.Map{})

}

package mongodb

import (
	"context"
	"testing"

	"github.com/mitchellh/mapstructure"
	mongoke "github.com/remorses/mongoke/src"
	"github.com/remorses/mongoke/src/testutil"
	"github.com/stretchr/testify/require"
)

func TestMongodbFunctions(t *testing.T) {
	collection := "users"
	ctx := context.Background()
	uri := testutil.MONGODB_URI
	type user struct {
		Name string `json:name`
		Age  int    `json:age`
	}

	type userStruct struct {
		Name string
		Age  int
	}

	exampleUsers := []mongoke.Map{
		{"name": "01", "age": 1},
		{"name": "02", "age": 2},
		{"name": "03", "age": 3},
	}
	// clear and insert some docs
	m := MongodbDatabaseFunctions{}
	db, err := m.InitMongo(uri)
	if err != nil {
		t.Error(err)
	}
	// clear
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

	t.Run("FindOne with eq", func(t *testing.T) {
		m := MongodbDatabaseFunctions{}
		users, err := m.FindMany(
			mongoke.FindManyParams{
				Collection:  collection,
				DatabaseUri: uri,
				Limit:       1,
				Where: map[string]mongoke.Filter{
					"name": {Eq: "01"},
				},
			},
		)
		user := users[0]
		if err != nil {
			t.Error(err)
		}
		t.Log(testutil.Pretty(user))
		var x userStruct
		if err := mapstructure.Decode(user, &x); err != nil {
			t.Error(err)
		}
		require.Equal(t, x.Name, "01")
		require.Equal(t, x.Age, 1)
	})
	t.Run("FindMany with first", func(t *testing.T) {
		m := MongodbDatabaseFunctions{}
		users, err := m.FindMany(
			mongoke.FindManyParams{
				Collection:  collection,
				DatabaseUri: uri,
				OrderBy:     map[string]int{"_id": mongoke.ASC},
				Limit:       2,
			},
		)
		if err != nil {
			t.Error(err)
		}
		t.Log(testutil.Pretty(users))
		var x []userStruct
		if err := mapstructure.Decode(users, &x); err != nil {
			t.Error(err)
		}
		require.Equal(t, 2, len(x))
		require.Equal(t, x[0].Name, "01")
		require.Equal(t, x[1].Name, "02")
	})
	t.Run("FindMany with neq", func(t *testing.T) {
		m := MongodbDatabaseFunctions{}
		users, err := m.FindMany(
			mongoke.FindManyParams{
				Collection:  collection,
				DatabaseUri: uri,
				Where: map[string]mongoke.Filter{
					"name": {Neq: "01"},
				},
				OrderBy: map[string]int{"_id": mongoke.ASC},
				Limit:   2,
			},
		)
		if err != nil {
			t.Error(err)
		}
		t.Log(testutil.Pretty(users))
		var x []userStruct
		if err := mapstructure.Decode(users, &x); err != nil {
			t.Error(err)
		}
		require.Equal(t, 2, len(x))
		require.Equal(t, x[0].Name, "02")
		require.Equal(t, x[1].Name, "03")
	})
	t.Run("FindMany direction DESC", func(t *testing.T) {
		m := MongodbDatabaseFunctions{}
		users, err := m.FindMany(
			mongoke.FindManyParams{
				Collection:  collection,
				DatabaseUri: uri,
				OrderBy:     map[string]int{"_id": mongoke.DESC},
				Limit:       2,
			},
		)
		if err != nil {
			t.Error(err)
		}
		t.Log(testutil.Pretty(users))
		var x []userStruct
		if err := mapstructure.Decode(users, &x); err != nil {
			t.Error(err)
		}
		require.Equal(t, 2, len(x))
		require.Equal(t, x[0].Name, "03")
		require.Equal(t, x[1].Name, "02")
	})
}

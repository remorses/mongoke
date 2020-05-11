package mongoke

import (
	"context"
	"testing"

	"github.com/mitchellh/mapstructure"
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

	exampleUsers := []Map{
		{"name": "01", "age": 1},
		{"name": "02", "age": 2},
		{"name": "03", "age": 3},
	}
	// clear and insert some docs
	m := MongodbDatabaseFunctions{}
	db, err := m.initMongo(uri)
	if err != nil {
		t.Error(err)
	}
	// clear
	_, err = db.Collection(collection).DeleteMany(ctx, Map{})
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
		user, err := m.FindOne(
			FindOneParams{
				Collection:  collection,
				DatabaseUri: uri,
				Where: map[string]Filter{
					"name": {Eq: "01"},
				},
			},
		)
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
			FindManyParams{
				Collection:  collection,
				DatabaseUri: uri,
				Direction:   ASC,
				Pagination: Pagination{
					First: 2,
				},
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
			FindManyParams{
				Collection:  collection,
				DatabaseUri: uri,
				Where: map[string]Filter{
					"name": {Neq: "01"},
				},
				Pagination: Pagination{
					First: 2,
				},
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
			FindManyParams{
				Collection:  collection,
				DatabaseUri: uri,
				Direction:   DESC,
				Pagination: Pagination{
					First: 2,
				},
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

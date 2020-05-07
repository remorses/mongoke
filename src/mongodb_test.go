package mongoke

import (
	"context"
	"testing"

	"github.com/remorses/mongoke/src/testutil"
	"go.mongodb.org/mongo-driver/bson"
)

func TestInitMongo(t *testing.T) {
	t.Run("init mongo", func(t *testing.T) {
		db, err := initMongo(testutil.MONGODB_URI)
		if err != nil {
			t.Error(err)
		}
		names, err := db.Client().ListDatabaseNames(context.TODO(), bson.D{{}})
		if err != nil {
			t.Error(err)
		}
		prettyPrint(names)
	})
	t.Run("findOne", func(t *testing.T) {
		db, err := initMongo(testutil.MONGODB_URI)
		if err != nil {
			t.Error(err)
		}
		x, err := findOne(db.Collection("users"), map[string]interface{}{"eq": "ciao"})
		if err != nil {
			t.Error(err)
		}
		prettyPrint("findOne", x)
	})
	t.Run("findMany", func(t *testing.T) {
		db, err := initMongo(testutil.MONGODB_URI)
		coll := db.Collection("users")
		coll.InsertMany(context.TODO(), []interface{}{bson.M{"name": "xxx"}, bson.M{"name": "yyy"}})
		if err != nil {
			t.Error(err)
		}
		x, err := findMany(
			coll,
			map[string]interface{}{"name": map[string]interface{}{"eq": "xxx"}},
			Pagination{first: 10},
			"name",
			ASC,
		)
		if err != nil {
			t.Error(err)
		}
		prettyPrint("findMany", x)
	})
	t.Run("findOne nil", func(t *testing.T) {
		db, err := initMongo(testutil.MONGODB_URI)
		if err != nil {
			t.Error(err)
		}
		x, err := findOne(db.Collection("users"), nil)
		if err != nil {
			t.Error(err)
		}
		prettyPrint("findOne nil", x)
	})
}

type Match = map[string]interface{}

func TestRewriteFilter(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		match := rewriteFilter(
			Match{"x": Match{"eq": "xxx"}},
		)
		prettyPrint("rewritten query", match)
	})
	t.Run("nil", func(t *testing.T) {
		match := rewriteFilter(
			nil,
		)
		prettyPrint("rewritten nil query", match)
	})
}

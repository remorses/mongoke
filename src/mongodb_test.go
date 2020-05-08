package mongoke

import (
	"context"
	"fmt"
	"testing"

	"github.com/remorses/mongoke/src/testutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestInitMongo(t *testing.T) {
	db, err := initMongo(testutil.MONGODB_URI)
	if err != nil {
		t.Error(err)
	}
	t.Run("init mongo", func(t *testing.T) {

		names, err := db.Client().ListDatabaseNames(context.TODO(), bson.D{{}})
		if err != nil {
			t.Error(err)
		}
		prettyPrint(names)
	})
	// t.Run("findOne", func(t *testing.T) {

	// 	x, err := findOne(FindOneParams{Collection: "users", D})
	// 	if err != nil {
	// 		t.Error(err)
	// 	}
	// 	prettyPrint("findOne", x)
	// })
	// t.Run("findMany", func(t *testing.T) {
	// 	coll := db.Collection("users")
	// 	coll.InsertMany(context.TODO(), []interface{}{bson.M{"name": "aaa", "obj": bson.M{"xxx": bson.M{"yyy": 3}}}})
	// 	if err != nil {
	// 		t.Error(err)
	// 	}
	// 	x, err := findMany(
	// 		coll,
	// 		map[string]interface{}{"name": map[string]interface{}{"eq": "aaa"}},
	// 		Pagination{First: 10},
	// 		"name",
	// 		ASC,
	// 	)
	// 	if err != nil {
	// 		t.Error(err)
	// 	}
	// 	prettyPrint("findMany", x)
	// })
	// t.Run("findOne nil", func(t *testing.T) {
	// 	x, err := findOne(db.Collection("users"), nil)
	// 	if err != nil {
	// 		t.Error(err)
	// 	}
	// 	prettyPrint("findOne nil", x)
	// })
	t.Run("findOne with Filter", func(t *testing.T) {
		// opts := options.FindOne().SetProjection(bson.M{"name": Filter{Eq: "xxx"}})
		// prettyPrint("opts", opts.Projection)
		// fmt.Println(opts.Projection)
		// res, err := bson.Marshal(bson.M{"name": Filter{Eq: "xxx"}})
		res := db.Collection("users").FindOne(context.TODO(), bson.M{"name": Filter{Eq: "aaa"}})
		if res.Err() == mongo.ErrNoDocuments {
			t.Log("no docs")
		}
		if res.Err() != nil {
			t.Error(err)
		}
		var x bson.M
		res.Decode(&x)
		prettyPrint("findOne with filter", x)
		fmt.Printf("%v\n", x)
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

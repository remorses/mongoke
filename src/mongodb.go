package mongoke

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

const TIMEOUT_CONNECT = 5

func initMongo(uri string) (*mongo.Database, error) {
	uriOptions, err := connstring.Parse(uri)
	if err != nil {
		return nil, err
	}
	dbName := uriOptions.Database
	if dbName == "" {
		return nil, errors.New("the db uri must contain the database name")
	}
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT_CONNECT*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return client.Database(dbName), nil
}

// type findOneParams struct {
// 	collection
// 	database
// }

func findOne(collection *mongo.Collection, _filter interface{}) (interface{}, error) {
	filter, ok := _filter.(map[string]interface{})
	if !ok && filter != nil {
		return nil, errors.New("the where argument filter must be an object or nil")
	}
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT_FIND*time.Second)
	res := collection.FindOne(ctx, rewriteFilter(filter))
	if res.Err() == mongo.ErrNoDocuments {
		return nil, nil
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	var document interface{}
	err := res.Decode(&document)
	if err != nil {
		return nil, err
	}
	return document, nil
}

// only works if the mongo operators are at second level of the match, like { field: { eq: "xxx" } }
func rewriteFilter(filter map[string]interface{}) map[string]map[string]interface{} {
	newFilter := make(map[string]map[string]interface{})
	for k, v := range filter {
		if v, ok := v.(map[string]interface{}); ok {
			newFilter[k] = addDollarSigns(v)
		}
	}
	return newFilter
}

func addDollarSigns(filter map[string]interface{}) map[string]interface{} {
	newFilter := make(map[string]interface{})
	for k, v := range filter {
		newFilter["$"+k] = v
	}
	return newFilter
}

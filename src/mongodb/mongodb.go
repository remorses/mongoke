package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	mongoke "github.com/remorses/mongoke/src"
	"github.com/remorses/mongoke/src/testutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

var (
	TIMEOUT_CONNECT = 5
	MAX_QUERY_TIME  = 10
	TIMEOUT_FIND    = 10
)

type MongodbDatabaseFunctions struct {
	db     *mongo.Database
	Config mongoke.Config
}

func (self MongodbDatabaseFunctions) databaseUri() string {
	return self.Config.Mongodb.Uri
}

func (self *MongodbDatabaseFunctions) FindMany(ctx context.Context, p mongoke.FindManyParams) ([]mongoke.Map, error) {
	db, err := self.Init(ctx)
	if err != nil {
		return nil, err
	}
	ctx, _ = context.WithTimeout(ctx, time.Duration(TIMEOUT_FIND)*time.Second)
	opts := options.Find()
	opts.SetMaxTime(time.Duration(MAX_QUERY_TIME) * time.Second)
	opts.SetLimit(int64(p.Limit))
	opts.SetSkip(int64(p.Offset))
	opts.SetSort(p.OrderBy)
	testutil.PrettyPrint(p)

	where := MakeMongodbMatch(p.Where)
	res, err := db.Collection(p.Collection).Find(ctx, where, opts)
	if err != nil {
		// log.Print("Error in findMany", err)
		return nil, err
	}
	defer res.Close(ctx)
	nodes := make([]mongoke.Map, 0)
	err = res.All(ctx, &nodes)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func (self *MongodbDatabaseFunctions) InsertMany(ctx context.Context, p mongoke.InsertManyParams) (mongoke.NodesMutationPayload, error) {
	payload := mongoke.NodesMutationPayload{}
	if len(p.Data) == 0 {
		return payload, nil
	}
	db, err := self.Init(ctx)
	if err != nil {
		return payload, err
	}
	opts := options.InsertMany()
	opts.SetOrdered(true)
	testutil.PrettyPrint(p)
	var data = make([]interface{}, len(p.Data))
	for i, x := range p.Data {
		data[i] = x
	}
	res, err := db.Collection(p.Collection).InsertMany(ctx, data, opts)
	if err != nil {
		// log.Print("Error in findMany", err)
		return payload, err
	}
	fmt.Println(res.InsertedIDs)
	for i, id := range res.InsertedIDs {
		p.Data[i]["_id"] = id
	}
	return mongoke.NodesMutationPayload{
		Returning:     p.Data[:len(res.InsertedIDs)],
		AffectedCount: len(res.InsertedIDs),
	}, nil
}

func (self *MongodbDatabaseFunctions) UpdateOne(ctx context.Context, p mongoke.UpdateParams) (mongoke.NodeMutationPayload, error) {
	db, err := self.Init(ctx)
	if err != nil {
		return mongoke.NodeMutationPayload{}, err
	}
	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(options.After)
	testutil.PrettyPrint(p)

	where := MakeMongodbMatch(p.Where)
	res := db.Collection(p.Collection).FindOneAndUpdate(ctx, where, bson.M{"$set": p.Set}, opts)
	if res.Err() == mongo.ErrNoDocuments {
		println("no docs to update")
		return mongoke.NodeMutationPayload{
			AffectedCount: 0,
			Returning:     nil,
		}, nil
	} else if res.Err() != nil {
		return mongoke.NodeMutationPayload{}, err
	}
	data := mongoke.Map{}
	err = res.Decode(&data)
	if err != nil {
		return mongoke.NodeMutationPayload{}, err
	}
	return mongoke.NodeMutationPayload{
		AffectedCount: 1,
		Returning:     data,
	}, nil
}

// first updateMany documents, then query again the documents and return them, all inside a transaction that prevents other writes happen before the query
func (self *MongodbDatabaseFunctions) UpdateMany(ctx context.Context, p mongoke.UpdateParams) (mongoke.NodesMutationPayload, error) {
	db, err := self.Init(ctx)
	payload := mongoke.NodesMutationPayload{}
	if err != nil {
		return payload, err
	}
	opts := options.Update()

	testutil.PrettyPrint(p)

	// TODO execute inside a transaction
	nodes, err := self.FindMany(ctx, mongoke.FindManyParams{Collection: p.Collection, Where: p.Where})
	if err != nil {
		return payload, err
	}

	where := MakeMongodbMatch(p.Where)
	res, err := db.Collection(p.Collection).UpdateMany(ctx, where, bson.M{"$set": p.Set}, opts)
	if err != nil {
		return payload, err
	}
	payload.AffectedCount = int(res.ModifiedCount + res.UpsertedCount)

	return mongoke.NodesMutationPayload{
		AffectedCount: payload.AffectedCount,
		Returning:     nodes,
	}, nil
}

func (self *MongodbDatabaseFunctions) DeleteMany(ctx context.Context, p mongoke.DeleteManyParams) (mongoke.NodesMutationPayload, error) {
	db, err := self.Init(ctx)
	payload := mongoke.NodesMutationPayload{}
	if err != nil {
		return payload, err
	}
	opts := options.Delete()

	testutil.PrettyPrint(p)

	nodes, err := self.FindMany(ctx, mongoke.FindManyParams{Collection: p.Collection, Where: p.Where})
	if err != nil {
		return payload, err
	}

	// TODO delete only documents user has permissions to

	where := MakeMongodbMatch(p.Where)
	res, err := db.Collection(p.Collection).DeleteMany(ctx, where, opts)
	if err != nil {
		return payload, err
	}
	return mongoke.NodesMutationPayload{
		AffectedCount: int(res.DeletedCount),
		Returning:     nodes,
	}, nil
}

func (self *MongodbDatabaseFunctions) Init(ctx context.Context) (*mongo.Database, error) {
	if self.db != nil {
		return self.db, nil
	}
	uri := self.databaseUri()
	if uri == "" {
		return nil, errors.New("uri is missing")
	}
	uriOptions, err := connstring.Parse(uri)
	if err != nil {
		return nil, err
	}
	dbName := uriOptions.Database
	if dbName == "" {
		return nil, errors.New("the db uri must contain the database name")
	}
	ctx, _ = context.WithTimeout(ctx, time.Duration(TIMEOUT_CONNECT)*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	db := client.Database(dbName)
	self.db = db
	return db, nil
}

func MakeMongodbMatch(where mongoke.WhereTree) map[string]interface{} {
	// TODO for every k, v use mapstructure to map to a filter
	// if k is or, and, use mapstructure to map to an array of filters
	var res = make(map[string]interface{})
	for k, v := range where.Match {
		res[k] = v
	}

	if len(where.And) != 0 {
		res["$and"] = make([]map[string]interface{}, len(where.And))
		for i, a := range where.And {
			res["$and"].([]map[string]interface{})[i] = MakeMongodbMatch(a)
		}
	}
	if len(where.Or) != 0 {
		res["$or"] = make([]map[string]interface{}, len(where.Or))
		for i, a := range where.Or {
			res["$or"].([]map[string]interface{})[i] = MakeMongodbMatch(a)
		}
	}
	return res
}

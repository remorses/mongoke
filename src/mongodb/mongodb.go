package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	goke "github.com/remorses/goke/src"
	"github.com/remorses/goke/src/testutil"
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
	Config goke.Config
}

func (self MongodbDatabaseFunctions) databaseUri() string {
	return self.Config.Mongodb.Uri
}

func (self *MongodbDatabaseFunctions) FindMany(ctx context.Context, p goke.FindManyParams, hook goke.TransformDocument) ([]goke.Map, error) {
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
	nodes := make([]goke.Map, 0)
	err = res.All(ctx, &nodes)
	if err != nil {
		return nil, err
	}
	nodes, err = goke.FilterDocuments(nodes, hook)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func (self *MongodbDatabaseFunctions) InsertMany(ctx context.Context, p goke.InsertManyParams, hook goke.TransformDocument) (goke.NodesMutationPayload, error) {
	payload := goke.NodesMutationPayload{}
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
	nodes, err := goke.FilterDocuments(p.Data, hook)
	if err != nil {
		return payload, err
	}
	res, err := db.Collection(p.Collection).InsertMany(ctx, goke.MapsToInterfaces(p.Data), opts)
	if err != nil {
		// log.Print("Error in findMany", err)
		return payload, err
	}
	fmt.Println(res.InsertedIDs)
	for i, id := range res.InsertedIDs {
		nodes[i]["_id"] = id
	}
	return goke.NodesMutationPayload{
		Returning:     nodes[:len(res.InsertedIDs)],
		AffectedCount: len(res.InsertedIDs),
	}, nil
}

func (self *MongodbDatabaseFunctions) UpdateOne(ctx context.Context, p goke.UpdateParams, hook goke.TransformDocument) (goke.NodeMutationPayload, error) {
	db, err := self.Init(ctx)
	payload := goke.NodeMutationPayload{
		Returning:     nil,
		AffectedCount: 0,
	}
	if err != nil {
		return payload, err
	}

	// TODO this step of checking could be skipped if there are no guards
	nodes, err := self.FindMany(ctx, goke.FindManyParams{
		Collection: p.Collection,
		Limit:      1,
		Where:      p.Where,
	}, hook)
	if err != nil {
		return payload, err
	}
	if len(nodes) == 0 {
		return payload, nil
	}

	// make sure we update the document we checked with the hook
	where := goke.ExtendWhereMatch(
		p.Where,
		map[string]goke.Filter{
			"_id": {
				Eq: nodes[0]["_id"],
			},
		},
	)
	match := MakeMongodbMatch(where)

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(options.After)
	testutil.PrettyPrint(p)

	res := db.Collection(p.Collection).FindOneAndUpdate(ctx, match, bson.M{"$set": p.Set}, opts)
	if res.Err() == mongo.ErrNoDocuments {
		println("no docs to update")
		return payload, nil
	} else if res.Err() != nil {
		return payload, err
	}
	data := goke.Map{}
	err = res.Decode(&data)
	if err != nil {
		return payload, err
	}
	return goke.NodeMutationPayload{
		AffectedCount: 1,
		Returning:     data,
	}, nil
}

// first updateMany documents, then query again the documents and return them, all inside a transaction that prevents other writes happen before the query
func (self *MongodbDatabaseFunctions) UpdateMany(ctx context.Context, p goke.UpdateParams, hook goke.TransformDocument) (goke.NodesMutationPayload, error) {
	payload := goke.NodesMutationPayload{}
	testutil.PrettyPrint(p)

	// TODO execute inside a transaction
	nodes, err := self.FindMany(ctx, goke.FindManyParams{Collection: p.Collection, Where: p.Where}, hook)
	if err != nil {
		return payload, err
	}

	for _, node := range nodes {
		r, err := self.UpdateOne(
			ctx,
			goke.UpdateParams{
				Collection: p.Collection,
				Set:        p.Set,
				Where: goke.ExtendWhereMatch(
					p.Where,
					map[string]goke.Filter{
						"_id": {
							Eq: node["_id"],
						},
					},
				),
			},
			nil, // pass nil to not repeat the check
		)
		if err != nil {
			return payload, err
		}
		payload.AffectedCount += r.AffectedCount
		if r.Returning != nil {
			payload.Returning = append(payload.Returning, r.Returning.(goke.Map))
		}
	}
	return payload, nil
}

func (self *MongodbDatabaseFunctions) DeleteMany(ctx context.Context, p goke.DeleteManyParams, hook goke.TransformDocument) (goke.NodesMutationPayload, error) {
	db, err := self.Init(ctx)
	payload := goke.NodesMutationPayload{}
	if err != nil {
		return payload, err
	}
	opts := options.Delete()

	testutil.PrettyPrint(p)

	nodes, err := self.FindMany(ctx, goke.FindManyParams{Collection: p.Collection, Where: p.Where}, hook)
	if err != nil {
		return payload, err
	}

	// TODO delete only documents user has permissions to

	where := MakeMongodbMatch(p.Where)
	res, err := db.Collection(p.Collection).DeleteMany(ctx, where, opts)
	if err != nil {
		return payload, err
	}
	return goke.NodesMutationPayload{
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

func MakeMongodbMatch(where goke.WhereTree) map[string]interface{} {
	// for every k, v use mapstructure to map to a filter
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

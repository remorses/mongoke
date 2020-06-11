package fakedata

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/256dpi/lungo"
	"github.com/pkg/errors"
	goke "github.com/remorses/goke/src"
	mongodb "github.com/remorses/goke/src/mongodb"
	"github.com/remorses/goke/src/testutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	TIMEOUT_CONNECT                  = 5
	MAX_QUERY_TIME                   = 10
	TIMEOUT_FIND                     = 10
	DEFAULT_DOCUMENTS_PER_COLLECTION = 50
)

type FakeDatabaseFunctions struct {
	Config             goke.Config
	skipDataGeneration bool
	db                 lungo.IDatabase
}

func (self *FakeDatabaseFunctions) FindMany(ctx context.Context, p goke.FindManyParams, hook goke.TransformDocument) ([]goke.Map, error) {
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

	where := mongodb.MakeMongodbMatch(p.Where)
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
	nodes, err = goke.FilterDocuments(nodes, hook) // TODO filtering after will potentially remove documents, making the limit lower than the query one
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func (self *FakeDatabaseFunctions) InsertMany(ctx context.Context, p goke.InsertManyParams, hook goke.TransformDocument) (goke.NodesMutationPayload, error) {
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

	// copy to new list of maps to not mutate input argument
	resultNodes := make([]goke.Map, len(res.InsertedIDs))
	for i, id := range res.InsertedIDs {
		if resultNodes[i] == nil {
			resultNodes[i] = make(goke.Map)
		}
		for k, v := range nodes[i] {
			resultNodes[i][k] = v
		}
		resultNodes[i]["_id"] = id
	}
	return goke.NodesMutationPayload{
		Returning:     resultNodes,
		AffectedCount: len(res.InsertedIDs),
	}, nil
}

// first updateMany documents, then query again the documents and return them, all inside a transaction that prevents other writes happen before the query
func (self *FakeDatabaseFunctions) UpdateMany(ctx context.Context, p goke.UpdateParams, hook goke.TransformDocument) (goke.NodesMutationPayload, error) {
	if p.Limit == 0 {
		p.Limit = math.MaxInt16
	}
	db, err := self.Init(ctx)
	payload := goke.NodesMutationPayload{
		Returning:     nil,
		AffectedCount: 0,
	}
	if err != nil {
		return payload, err
	}
	testutil.PrettyPrint(p)

	// TODO execute inside a transaction
	nodes, err := self.FindMany(ctx, goke.FindManyParams{Collection: p.Collection, Where: p.Where, Limit: p.Limit}, hook)
	if err != nil {
		return payload, err
	}

	for _, node := range nodes {
		match := mongodb.MakeMongodbMatch(goke.ExtendWhereMatch(
			p.Where, // TODO maybe do not extend, if a concurrent update makes the document no more suitable we lose the update, maybe better?
			map[string]goke.Filter{
				"_id": {
					Eq: node["_id"],
				},
			},
		))
		opts := options.FindOneAndUpdate()
		opts.SetReturnDocument(options.After)
		testutil.PrettyPrint(p)

		res := db.Collection(p.Collection).FindOneAndUpdate(ctx, match, bson.M{"$set": p.Set}, opts)
		if res.Err() == mongo.ErrNoDocuments {
			return payload, nil
		} else if res.Err() != nil {
			return payload, err
		}
		data := goke.Map{}
		err = res.Decode(&data)
		if err != nil {
			return payload, err
		}
		payload.AffectedCount++
		payload.Returning = append(payload.Returning, data)
	}
	return payload, nil
}

func (self *FakeDatabaseFunctions) DeleteMany(ctx context.Context, p goke.DeleteManyParams, hook goke.TransformDocument) (goke.NodesMutationPayload, error) {
	if p.Limit == 0 {
		p.Limit = math.MaxInt16
	}
	db, err := self.Init(ctx)
	payload := goke.NodesMutationPayload{}
	if err != nil {
		return payload, err
	}
	opts := options.Delete()
	testutil.PrettyPrint(p)

	// find accessible nodes
	nodes, err := self.FindMany(ctx, goke.FindManyParams{Collection: p.Collection, Where: p.Where, Limit: p.Limit}, hook)
	if err != nil {
		return payload, err
	}

	for _, node := range nodes {
		// find one specific node
		where := goke.ExtendWhereMatch(
			p.Where,
			map[string]goke.Filter{
				"_id": {
					Eq: node["_id"],
				},
			},
		)
		match := mongodb.MakeMongodbMatch(where)
		res, err := db.Collection(p.Collection).DeleteOne(ctx, match, opts)
		if err != nil {
			return payload, err
		}
		payload.AffectedCount += int(res.DeletedCount)
		payload.Returning = append(payload.Returning, node)
	}

	return payload, nil
}

func (self *FakeDatabaseFunctions) Init(ctx context.Context) (lungo.IDatabase, error) {
	if self.db != nil {
		return self.db, nil
	}
	println("initializing fake data client")
	opts := lungo.Options{
		Store: lungo.NewMemoryStore(),
	}
	client, _, err := lungo.Open(nil, opts)
	if err != nil {
		return nil, err
	}

	// ensure engine is closed
	// defer engine.Close()

	// get db
	db := client.Database("default")
	self.db = db

	if self.skipDataGeneration {
		return db, nil
	}

	err = self.generateFakeData(self.Config)
	if err != nil {
		return nil, errors.Wrap(err, "Error generating fake data")
	}
	return db, nil
}

func (self FakeDatabaseFunctions) generateFakeData(config goke.Config) error {
	println("generating fake data")
	faker, err := NewFakeData(NewFakeDataParams{typeDefs: config.Schema})
	if err != nil {
		return err
	}
	documentsPerCollection := config.FakeDatabase.DocumentsPerCollection
	if documentsPerCollection == nil {
		documentsPerCollection = &DEFAULT_DOCUMENTS_PER_COLLECTION
	}
	if *documentsPerCollection == 0 {
		return nil
	}
	for name, t := range config.Types {
		var docs []interface{}

		for i := 0; i < *documentsPerCollection; i++ {
			data, err := faker.Generate(name)
			if err != nil {
				return err
			}
			docs = append(docs, data)
		}
		self.db.Collection(t.Collection).InsertMany(context.Background(), docs)
	}
	return nil
}

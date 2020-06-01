package fakedata

import (
	"context"
	"fmt"
	"time"

	"github.com/256dpi/lungo"
	"github.com/pkg/errors"
	mongoke "github.com/remorses/mongoke/src"
	"github.com/remorses/mongoke/src/testutil"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	TIMEOUT_CONNECT                  = 5
	MAX_QUERY_TIME                   = 10
	TIMEOUT_FIND                     = 10
	DEFAULT_DOCUMENTS_PER_COLLECTION = 50
)

type FakeDatabaseFunctions struct {
	mongoke.DatabaseInterface
	Config             mongoke.Config
	skipDataGeneration bool
	db                 lungo.IDatabase
}

func (self *FakeDatabaseFunctions) FindMany(ctx context.Context, p mongoke.FindManyParams) ([]mongoke.Map, error) {
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

	res, err := db.Collection(p.Collection).Find(ctx, p.Where, opts)
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

func (self *FakeDatabaseFunctions) InsertMany(ctx context.Context, p mongoke.InsertManyParams) ([]mongoke.Map, error) {
	db, err := self.Init(ctx)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	fmt.Println(res.InsertedIDs)
	for i, id := range res.InsertedIDs {
		p.Data[i]["_id"] = id
	}
	return p.Data, nil
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

func (self FakeDatabaseFunctions) generateFakeData(config mongoke.Config) error {
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

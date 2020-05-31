package fakedata

import (
	"context"
	"time"

	"github.com/256dpi/lungo"
	mongoke "github.com/remorses/mongoke/src"
	"github.com/remorses/mongoke/src/testutil"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	TIMEOUT_CONNECT = 5
	MAX_QUERY_TIME  = 10
	TIMEOUT_FIND    = 10
)

type FakeDatabaseFunctions struct {
	mongoke.DatabaseInterface
	db lungo.IDatabase
}

func (self FakeDatabaseFunctions) FindMany(ctx context.Context, p mongoke.FindManyParams) ([]mongoke.Map, error) {
	db, err := self.Init(ctx, p.DatabaseUri)
	if err != nil {
		return nil, err
	}
	ctx, _ = context.WithTimeout(ctx, TIMEOUT_FIND*time.Second)
	opts := options.Find()
	opts.SetMaxTime(MAX_QUERY_TIME * time.Second)
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

func (self *FakeDatabaseFunctions) Init(ctx context.Context, uri string) (lungo.IDatabase, error) {
	if self.db != nil {
		return self.db, nil
	}
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
	return db, nil
}

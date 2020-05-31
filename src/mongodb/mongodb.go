package mongodb

import (
	"context"
	"errors"
	"time"

	mongoke "github.com/remorses/mongoke/src"
	"github.com/remorses/mongoke/src/testutil"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

const (
	TIMEOUT_CONNECT = 5
	MAX_QUERY_TIME  = 10
	TIMEOUT_FIND    = 10
)

type MongodbDatabaseFunctions struct {
	mongoke.DatabaseInterface
	db *mongo.Database
}

func (self *MongodbDatabaseFunctions) FindMany(ctx context.Context, p mongoke.FindManyParams) ([]mongoke.Map, error) {
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

func (self *MongodbDatabaseFunctions) Init(ctx context.Context, uri string) (*mongo.Database, error) {
	if self.db != nil {
		return self.db, nil
	}
	uriOptions, err := connstring.Parse(uri)
	if err != nil {
		return nil, err
	}
	dbName := uriOptions.Database
	if dbName == "" {
		return nil, errors.New("the db uri must contain the database name")
	}
	ctx, _ = context.WithTimeout(ctx, TIMEOUT_CONNECT*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	db := client.Database(dbName)
	self.db = db
	return db, nil
}

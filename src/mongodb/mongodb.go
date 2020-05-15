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

// type findOneParams struct {
// 	collection
// 	database
// }

type MongodbDatabaseFunctions struct {
	mongoke.DatabaseInterface
	db *mongo.Database
}

// func (self MongodbDatabaseFunctions) FindOne(p mongoke.FindOneParams) (interface{}, error) {
// 	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT_FIND*time.Second)
// 	db, err := self.InitMongo(p.DatabaseUri)
// 	if err != nil {
// 		return nil, err
// 	}
// 	collection := db.Collection(p.Collection)
// 	// testutil.PrettyPrint(p.Where)

// 	opts := options.FindOne().SetMaxTime(MAX_QUERY_TIME * time.Second)
// 	res := collection.FindOne(ctx, p.Where, opts)

// 	if res.Err() == mongo.ErrNoDocuments {
// 		return nil, nil
// 	}
// 	if res.Err() != nil {
// 		return nil, res.Err()
// 	}
// 	var document mongoke.Map = make(mongoke.Map)
// 	err = res.Decode(document)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// testutil.PrettyPrint(document)
// 	return document, nil
// }

const (
	DEFAULT_NODES_COUNT = 20
	MAX_NODES_COUNT     = 40
)

func (self MongodbDatabaseFunctions) FindMany(p mongoke.FindManyParams) ([]mongoke.Map, error) {
	if p.Direction == 0 {
		p.Direction = mongoke.ASC
	}
	if p.CursorField == "" {
		p.CursorField = "_id"
	}
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT_FIND*time.Second)
	db, err := self.InitMongo(p.DatabaseUri)
	if err != nil {
		return nil, err
	}
	after := p.Pagination.After
	before := p.Pagination.Before
	last := p.Pagination.Last
	first := p.Pagination.First

	opts := options.Find()

	// set defaults
	if first == 0 && last == 0 {
		if after != "" {
			first = DEFAULT_NODES_COUNT
		} else if before != "" {
			last = DEFAULT_NODES_COUNT
		} else {
			first = DEFAULT_NODES_COUNT
		}
	}

	// TODO move findMany argument handling logic in custom function
	// assertion for arguments
	if after != "" && first == 0 && before == "" {
		return nil, errors.New("need `first` or `before` if using `after`")
	}
	if before != "" && (last == 0 && after == "") {
		return nil, errors.New("need `last` or `after` if using `before`")
	}
	if first != 0 && last != 0 {
		return nil, errors.New("cannot use `first` and `last` together")
	}

	// gt and lt
	cursorFieldMatch := p.Where[p.CursorField]
	if after != "" {
		if p.Direction == mongoke.DESC {
			cursorFieldMatch.Lt = after
		} else {
			cursorFieldMatch.Gt = after
		}
	}
	if before != "" {
		if p.Direction == mongoke.DESC {
			cursorFieldMatch.Gt = before
		} else {
			cursorFieldMatch.Lt = before
		}
	}

	// sort order
	sorting := p.Direction
	if last != 0 {
		sorting = -p.Direction
	}
	opts.SetSort(mongoke.Map{p.CursorField: sorting})

	// limit
	if last != 0 {
		opts.SetLimit(int64(min(MAX_NODES_COUNT, last)))
	}
	if first != 0 {
		opts.SetLimit(int64(min(MAX_NODES_COUNT, first)))
	}
	if first == 0 && last == 0 { // when using `after` and `before`
		opts.SetLimit(int64(MAX_NODES_COUNT))
	}

	opts.SetMaxTime(MAX_QUERY_TIME * time.Second)

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

func (self *MongodbDatabaseFunctions) InitMongo(uri string) (*mongo.Database, error) {
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
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT_CONNECT*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	db := client.Database(dbName)
	self.db = db
	return db, nil
}

// removes last or first node, adds pageInfo data

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

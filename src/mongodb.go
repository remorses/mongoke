package mongoke

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

const TIMEOUT_CONNECT = 5

// type findOneParams struct {
// 	collection
// 	database
// }

type MongodbDatabaseFunctions struct {
	db *mongo.Database
}

func (c MongodbDatabaseFunctions) FindOne(p FindOneParams) (interface{}, error) {
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT_FIND*time.Second)
	db, err := c.initMongo(p.DatabaseUri)
	if err != nil {
		return nil, err
	}
	collection := db.Collection(p.Collection)
	prettyPrint(p.Where)
	res := collection.FindOne(ctx, p.Where)

	if res.Err() == mongo.ErrNoDocuments {
		return nil, nil
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	var document bson.M = make(bson.M)
	err = res.Decode(document)
	if err != nil {
		return nil, err
	}
	prettyPrint(document)
	return document, nil
}

const (
	DEFAULT_NODES_COUNT = 20
	MAX_NODES_COUNT     = 40
)

const (
	ASC  = 1
	DESC = -1
)

func (c MongodbDatabaseFunctions) FindMany(p FindManyParams) (Connection, error) {
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT_FIND*time.Second)
	db, err := c.initMongo(p.DatabaseUri)
	if err != nil {
		return Connection{}, err
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

	// assertion for arguments
	if after != "" && (first == 0 || before == "") {
		return Connection{}, errors.New("need `first` or `before` if using `after`")
	}
	if before != "" && (last == 0 || after == "") {
		return Connection{}, errors.New("need `last` or `after` if using `before`")
	}
	if first != 0 && last != 0 {
		return Connection{}, errors.New("need `last` or `after` if using `before`")
	}

	// gt and lt
	cursorFieldMatch := p.Where[p.CursorField]
	if after != "" {
		if p.Direction == DESC {
			cursorFieldMatch.Lt = after
		} else {
			cursorFieldMatch.Gt = after
		}
	}
	if before != "" {
		if p.Direction == DESC {
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
	opts.SetSort(bson.M{p.CursorField: sorting})

	// limit
	if last != 0 {
		opts.SetLimit(int64(min(MAX_NODES_COUNT, last+1)))
	}
	if first != 0 {
		opts.SetLimit(int64(min(MAX_NODES_COUNT, first+1)))
	}
	if first == 0 && last == 0 { // when using `after` and `before`
		opts.SetLimit(int64(MAX_NODES_COUNT))
	}

	prettyPrint(p)

	res, err := db.Collection(p.Collection).Find(ctx, p.Where, opts)
	if err != nil {
		// log.Print("Error in findMany", err)
		return Connection{}, err
	}
	defer res.Close(ctx)
	nodes := make([]bson.M, 0)
	err = res.All(ctx, &nodes)
	if err != nil {
		return Connection{}, err
	}

	connection := makeConnection(nodes, p.Pagination, p.CursorField)

	return connection, nil
}

func (c *MongodbDatabaseFunctions) initMongo(uri string) (*mongo.Database, error) {
	if c.db != nil {
		return c.db, nil
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
	c.db = db
	return db, nil
}

// removes last or first node, adds pageInfo data

func makeConnection(nodes []bson.M, pagination Pagination, cursorField string) Connection {
	if len(nodes) == 0 {
		return Connection{}
	}
	var hasNext bool
	var hasPrev bool
	var endCursor interface{}
	var startCursor interface{}
	if pagination.First != 0 {
		hasNext = len(nodes) == int(pagination.First+1)
		if hasNext {
			nodes = nodes[:len(nodes)-1]
		}
	}
	if pagination.Last != 0 {
		nodes = reverse(nodes)
		hasPrev = len(nodes) == int(pagination.Last+1)
		if hasPrev {
			nodes = nodes[1:]
		}
	}
	if len(nodes) != 0 {
		endCursor = nodes[len(nodes)-1][cursorField]
		startCursor = nodes[0][cursorField]
	}
	return Connection{
		Nodes: nodes,
		Edges: makeEdges(nodes, cursorField),
		PageInfo: PageInfo{
			StartCursor:     startCursor,
			EndCursor:       endCursor,
			HasNextPage:     hasNext,
			HasPreviousPage: hasPrev,
		},
	}
}

func makeEdges(nodes []bson.M, cursorField string) []Edge {
	var edges []Edge
	for _, node := range nodes {
		edges = append(edges, Edge{
			Node:   node,
			Cursor: node[cursorField],
		})
	}
	return edges
}

func reverse(input []bson.M) []bson.M {
	if len(input) == 0 {
		return input
	}
	// TODO remove recursion
	return append(reverse(input[1:]), input[0])
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

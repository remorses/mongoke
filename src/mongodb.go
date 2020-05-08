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

func initMongo(uri string) (*mongo.Database, error) { // TODO dont reconnect every time, save the instance
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

type Filter struct {
	Eq  interface{}   `bson:"$eq,omitempty"`
	Neq interface{}   `bson:"$ne,omitempty"`
	In  []interface{} `bson:"$in,omitempty"`
	Nin []interface{} `bson:"$nin,omitempty"`
	Gt  interface{}   `bson:"$gt,omitempty"`
	Lt  interface{}   `bson:"$lt,omitempty"`
}

type FindOneParams struct {
	Collection  string
	DatabaseUri string
	Where       map[string]Filter `mapstructure:"where"`
}

func findOne(p FindOneParams) (interface{}, error) {
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT_FIND*time.Second)
	db, err := initMongo(p.DatabaseUri)
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

type Pagination struct {
	First  int    `mapstructure:first`
	Last   int    `mapstructure:last`
	After  string `mapstructure:after`
	Before string `mapstructure:before`
}

const (
	DEFAULT_NODES_COUNT = 40
)

const (
	ASC  = 1
	DESC = -1
)

type FindManyParams struct {
	Collection  string
	DatabaseUri string
	Where       map[string]Filter `mapstructure:"where"`
	Pagination  Pagination
	CursorField string `mapstructure:"cursorField"`
	Direction   int    `mapstructure:"direction"`
}

func findMany(p FindManyParams) (interface{}, error) {
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT_FIND*time.Second)
	db, err := initMongo(p.DatabaseUri)
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

	// assertion for arguments
	if after != "" && (first == 0 || before == "") {
		return nil, errors.New("need `first` or `before` if using `after`")
	}
	if before != "" && (last == 0 || after == "") {
		return nil, errors.New("need `last` or `after` if using `before`")
	}
	if first != 0 && last != 0 {
		return nil, errors.New("need `last` or `after` if using `before`")
	}

	// gt and lt
	cursorFieldMatch := p.Where[p.CursorField] // TODO add already existing match
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
		opts.SetLimit(int64(last + 1))
	}
	if first != 0 {
		opts.SetLimit(int64(first + 1))
	}

	// execute
	// prettyPrint(rewriteFilter(p.Where))
	res, err := db.Collection(p.Collection).Find(ctx, p.Where, opts)
	if err != nil {
		// log.Print("Error in findMany", err)
		return nil, err
	}
	defer res.Close(ctx)
	nodes := make([]bson.M, 0)
	err = res.All(ctx, &nodes) // TODO limit to maxlen
	if err != nil {
		return nil, err
	}

	connection := makeConnection(nodes, p.Pagination, p.CursorField)

	return connection, nil

}

type PageInfo struct {
	StartCursor     interface{} `json:startCursor` // TODO should be string hash made from any object (scalar, enum, ...)
	EndCursor       interface{} `json:endCursor`
	HasNextPage     bool        `json:hasNextPage`
	HasPreviousPage bool        `json:hasPreviousPage`
}

type Connection struct {
	Nodes    []bson.M `json:nodes`
	PageInfo PageInfo `json:pageInfo`
}

// removes last or first node, adds pageInfo data
func makeConnection(nodes []bson.M, pagination Pagination, cursorField string) Connection {
	if len(nodes) == 0 {
		return Connection{}
	}
	var hasNext bool
	var hasPrev bool
	var endCursor interface{} // TODO should be string
	var startCursor interface{}
	if pagination.First != 0 {
		hasNext = len(nodes) == int(pagination.First+1)
		if hasNext {
			nodes = nodes[:len(nodes)-1] // TODO is right?
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
		PageInfo: PageInfo{
			StartCursor:     startCursor,
			EndCursor:       endCursor,
			HasNextPage:     hasNext,
			HasPreviousPage: hasPrev,
		},
	}
}

// TODO use a Match type with eq, lt, ...
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

func reverse(input []bson.M) []bson.M {
	if len(input) == 0 {
		return input
	}
	// TODO remove recursion
	return append(reverse(input[1:]), input[0])
}

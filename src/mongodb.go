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
	Eq  interface{} `bson:"$eq"`
	neq interface{} `bson:"$neq"`
	in  interface{} `bson:"$in"`
	nin interface{} `bson:"$nin"`
}

type FindOneParams struct {
	Collection  string
	DatabaseUri string
	Where       map[string]Filter
}

func findOne(p FindOneParams) (interface{}, error) {
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT_FIND*time.Second)
	db, err := initMongo(p.DatabaseUri)
	if err != nil {
		return nil, err
	}
	collection := db.Collection(p.Collection)
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

func findMany(collection *mongo.Collection, _filter interface{}, pagination Pagination, cursorField string, direction int) (interface{}, error) {
	filter, ok := _filter.(map[string]interface{}) // TODO for testing it would be cooler to use bson.M
	if !ok && _filter != nil {
		return nil, errors.New("the where argument filter must be an object or nil")
	}
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT_FIND*time.Second)

	after := pagination.After
	before := pagination.Before
	last := pagination.Last
	first := pagination.First

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

	// set right lt and gt
	lt := "$lt"
	gt := "$gt"
	if direction == DESC {
		lt = "$gt"
		gt = "$lt"
	}

	// gt and lt
	cursorFieldMatch := make(map[string]interface{}) // TODO add already existing match
	if after != "" {
		cursorFieldMatch[gt] = after
		filter[cursorField] = cursorFieldMatch
	}
	if before != "" {
		cursorFieldMatch[lt] = before
		filter[cursorField] = cursorFieldMatch
	}

	// sort order
	sorting := direction
	if last != 0 {
		sorting = -direction
	}
	opts.SetSort(bson.M{cursorField: sorting})

	// limit
	if last != 0 {
		opts.SetLimit(int64(last + 1))
	}
	if first != 0 {
		opts.SetLimit(int64(first + 1))
	}

	// execute
	prettyPrint(rewriteFilter(filter))
	res, err := collection.Find(ctx, rewriteFilter(filter), opts)
	if err != nil {
		// log.Print("Error in findMany", err)
		return nil, err
	}
	defer res.Close(ctx)
	nodes := make([]map[string]interface{}, 0)
	err = res.All(ctx, &nodes) // TODO limit to maxlen
	if err != nil {
		return nil, err
	}

	connection := makeConnection(nodes, pagination, cursorField)

	return connection, nil

}

type PageInfo struct {
	StartCursor     interface{} `json:startCursor` // TODO should be string hash made from any object (scalar, enum, ...)
	EndCursor       interface{} `json:endCursor`
	HasNextPage     bool        `json:hasNextPage`
	HasPreviousPage bool        `json:hasPreviousPage`
}

type Connection struct {
	Nodes    []map[string]interface{} `json:nodes`
	PageInfo PageInfo                 `json:pageInfo`
}

// removes last or first node, adds pageInfo data
func makeConnection(nodes []map[string]interface{}, pagination Pagination, cursorField string) Connection {
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

func reverse(input []map[string]interface{}) []map[string]interface{} {
	if len(input) == 0 {
		return input
	}
	// TODO remove recursion
	return append(reverse(input[1:]), input[0])
}

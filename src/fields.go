package mongoke

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
)

const TIMEOUT_FIND = 10

type findOneFieldConfig struct {
	collection string
	returnType *graphql.Object
}

func (mongoke *Mongoke) findOneField(conf findOneFieldConfig) *graphql.Field {
	resolver := func(params graphql.ResolveParams) (interface{}, error) {
		args := params.Args
		opts := FindOneParams{
			Collection:  conf.collection,
			DatabaseUri: mongoke.databaseUri,
		}
		err := mapstructure.Decode(args, &opts)
		if err != nil {
			return nil, err
		}
		document, err := mongoke.databaseFunctions.FindOne(opts)
		if err != nil {
			return nil, err
		}
		// document, err := mongoke.database.findOne()
		// prettyPrint(args)
		return document, nil
	}
	whereArg, err := mongoke.getWhereArg(conf.returnType)
	if err != nil {
		panic(err)
	}
	return &graphql.Field{
		Type: conf.returnType,
		Args: graphql.FieldConfigArgument{
			"where": &graphql.ArgumentConfig{Type: whereArg},
		},
		Resolve: resolver,
	}
}

type findManyFieldConfig struct {
	collection string
	returnType *graphql.Object
}

func (mongoke *Mongoke) findManyField(conf findManyFieldConfig) *graphql.Field {
	resolver := func(params graphql.ResolveParams) (interface{}, error) {
		args := params.Args
		opts := FindManyParams{
			DatabaseUri: mongoke.databaseUri, // here i set the defaults
			Collection:  conf.collection,
			Direction:   ASC,
			CursorField: "_id",
			Pagination:  paginationFromArgs(args),
		}
		err := mapstructure.Decode(args, &opts)
		if err != nil {
			return nil, err
		}
		nodes, err := mongoke.databaseFunctions.FindMany(
			opts,
		)
		connection := makeConnection(nodes, opts.Pagination, opts.CursorField)
		if err != nil {
			return nil, err
		}
		// document, err := mongoke.database.findOne()
		// prettyPrint(args)
		return connection, nil
	}
	whereArg, err := mongoke.getWhereArg(conf.returnType)
	if err != nil {
		panic(err)
	}
	connectionType, err := mongoke.getConnectionType(conf.returnType)
	if err != nil {
		panic(err)
	}
	return &graphql.Field{
		Type: connectionType,
		Args: graphql.FieldConfigArgument{
			"where":       &graphql.ArgumentConfig{Type: whereArg},
			"first":       &graphql.ArgumentConfig{Type: graphql.Int},
			"last":        &graphql.ArgumentConfig{Type: graphql.Int},
			"after":       &graphql.ArgumentConfig{Type: AnyScalar},
			"before":      &graphql.ArgumentConfig{Type: AnyScalar},
			"direction":   &graphql.ArgumentConfig{Type: directionEnum},
			"cursorField": &graphql.ArgumentConfig{Type: graphql.String}, // TODO make cursorField as the indexable fields enum, so people dont get access to private fields
		},
		Resolve: resolver,
	}
}

func paginationFromArgs(args interface{}) Pagination {
	var pag Pagination
	err := mapstructure.Decode(args, &pag)
	if err != nil {
		fmt.Println(err)
		return Pagination{}
	}
	// prettyPrint(pag)
	return pag
}

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

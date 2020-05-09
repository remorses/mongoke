package mongoke

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
)

const TIMEOUT_FIND = 10

type createFieldParams struct {
	collection  string
	permissions []AuthGuard
	returnType  *graphql.Object
}

func (mongoke *Mongoke) findOneField(conf createFieldParams) *graphql.Field {
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
		jwt := Map{} // TODO take jwt from rootObject
		// don't compute permissions if document is nil
		document, err = applyGuardsOnDocument(applyGuardsOnDocumentParams{
			document:  document,
			guards:    conf.permissions,
			jwt:       jwt,
			operation: Operations.READ,
		})
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

func (mongoke *Mongoke) findManyField(conf createFieldParams) *graphql.Field {
	resolver := func(params graphql.ResolveParams) (interface{}, error) {
		args := params.Args
		pagination := paginationFromArgs(args)
		opts := FindManyParams{
			DatabaseUri: mongoke.databaseUri, // here i set the defaults
			Collection:  conf.collection,
			Direction:   ASC,
			CursorField: "_id",
			Pagination:  pagination,
		}
		err := mapstructure.Decode(args, &opts)
		if err != nil {
			return nil, err
		}
		nodes, err := mongoke.databaseFunctions.FindMany(
			opts,
		)
		if err != nil {
			return nil, err
		}

		// TODO skip auth if no auth guards present
		// if len(conf.permissions) == 0 {
		// 	return connection, nil
		// }

		jwt := Map{} // TODO take jwt from rootObject
		var accessibleNodes []Map
		for _, document := range nodes {
			node, err := applyGuardsOnDocument(applyGuardsOnDocumentParams{
				document:  document,
				guards:    conf.permissions,
				jwt:       jwt,
				operation: Operations.READ,
			})
			if err != nil {
				// println("got an error while calling applyGuardsOnDocument on findManyField for " + conf.returnType.PrivateName)
				// fmt.Println(err)
				continue
			}
			if node != nil {
				accessibleNodes = append(accessibleNodes, node.(Map))
			}
		}
		connection := makeConnection(accessibleNodes, opts.Pagination, opts.CursorField)
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
	indexableFieldsEnum, err := mongoke.getIndexableFieldsEnum(conf.returnType)
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
			"cursorField": &graphql.ArgumentConfig{Type: indexableFieldsEnum},
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
	// increment nodes count so createConnection knows how to set `hasNextPage`
	if pag.First != 0 {
		pag.First++
	}
	if pag.Last != 0 {
		pag.Last++
	}
	// prettyPrint(pag)
	return pag
}

func makeConnection(nodes []Map, pagination Pagination, cursorField string) Connection {
	if len(nodes) == 0 {
		return Connection{}
	}
	var hasNext bool
	var hasPrev bool
	var endCursor interface{}
	var startCursor interface{}
	if pagination.First != 0 {
		hasNext = len(nodes) == int(pagination.First)
		if hasNext {
			nodes = nodes[:len(nodes)-1]
		}
	}
	if pagination.Last != 0 {
		nodes = reverse(nodes)
		hasPrev = len(nodes) == int(pagination.Last)
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

func makeEdges(nodes []Map, cursorField string) []Edge {
	var edges []Edge
	for _, node := range nodes {
		edges = append(edges, Edge{
			Node:   node,
			Cursor: node[cursorField],
		})
	}
	return edges
}

func reverse(input []Map) []Map {
	if len(input) == 0 {
		return input
	}
	// TODO remove recursion
	return append(reverse(input[1:]), input[0])
}

func reverseStrings(input []string) []string {
	if len(input) == 0 {
		return input
	}
	// TODO remove recursion
	return append(reverseStrings(input[1:]), input[0])
}

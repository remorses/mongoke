package mongoke

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/graphql-go/graphql"
	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
)

const TIMEOUT_FIND = 10

type createFieldParams struct {
	collection   string
	initialWhere map[string]Filter
	permissions  []AuthGuard
	returnType   graphql.Type
	omitWhere    bool
}

func (mongoke *Mongoke) findOneField(conf createFieldParams) (*graphql.Field, error) {
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
		if conf.initialWhere != nil {
			mergo.Merge(&opts.Where, conf.initialWhere)
		}
		document, err := mongoke.databaseFunctions.FindOne(opts)
		if err != nil {
			return nil, err
		}
		jwt := getJwt(params)
		// don't compute permissions if document is nil
		if document == nil {
			return nil, nil
		}
		document, err = applyGuardsOnDocument(applyGuardsOnDocumentParams{
			document:  document,
			guards:    conf.permissions,
			jwt:       jwt,
			operation: Operations.READ,
		})
		if err != nil {
			return nil, err
		}
		return document, nil
	}
	whereArg, err := mongoke.getWhereArg(conf.returnType)
	if err != nil {
		return nil, err
	}
	args := graphql.FieldConfigArgument{}
	if !conf.omitWhere {
		args["where"] = &graphql.ArgumentConfig{Type: whereArg}
	}
	field := graphql.Field{
		Type:    conf.returnType,
		Args:    args,
		Resolve: resolver,
	}
	return &field, nil
}

func (mongoke *Mongoke) findManyField(conf createFieldParams) (*graphql.Field, error) {
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
		if conf.initialWhere != nil {
			mergo.Merge(&opts.Where, conf.initialWhere)
		}
		nodes, err := mongoke.databaseFunctions.FindMany(
			opts,
		)
		if err != nil {
			return nil, err
		}

		if len(conf.permissions) == 0 {
			return makeConnection(nodes, opts.Pagination, opts.CursorField), nil
		}

		jwt := getJwt(params)
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
		return nil, err
	}
	connectionType, err := mongoke.getConnectionType(conf.returnType)
	if err != nil {
		return nil, err
	}
	indexableFieldsEnum, err := mongoke.getIndexableFieldsEnum(conf.returnType)
	if err != nil {
		return nil, err
	}
	args := graphql.FieldConfigArgument{
		"first":       &graphql.ArgumentConfig{Type: graphql.Int},
		"last":        &graphql.ArgumentConfig{Type: graphql.Int},
		"after":       &graphql.ArgumentConfig{Type: AnyScalar},
		"before":      &graphql.ArgumentConfig{Type: AnyScalar},
		"direction":   &graphql.ArgumentConfig{Type: directionEnum},
		"cursorField": &graphql.ArgumentConfig{Type: indexableFieldsEnum},
	}
	if !conf.omitWhere {
		args["where"] = &graphql.ArgumentConfig{Type: whereArg}
	}
	field := graphql.Field{
		Type:    connectionType,
		Args:    args,
		Resolve: resolver,
	}
	return &field, nil
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

func getJwt(params graphql.ResolveParams) jwt.MapClaims {
	root := params.Info.RootValue
	rootMap, ok := root.(Map)
	if !ok {
		return jwt.MapClaims{}

	}
	v, ok := rootMap["jwt"]
	if !ok {
		return jwt.MapClaims{}
	}
	jwtMap, ok := v.(jwt.MapClaims)
	if !ok {
		return jwt.MapClaims{}
	}
	return jwtMap
}

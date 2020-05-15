package fields

import (
	"errors"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/graphql-go/graphql"
	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
	mongoke "github.com/remorses/mongoke/src"
	"github.com/remorses/mongoke/src/types"
)

type typeNodesArgs struct {
	Where       map[string]mongoke.Filter `mapstructure:"where"`
	Pagination  mongoke.Pagination
	CursorField string `mapstructure:"cursorField"`
	Direction   int    `mapstructure:"direction"`
}

type pageInfo struct {
	StartCursor     interface{} `json:"startCursor"`
	EndCursor       interface{} `json:"endCursor"`
	HasNextPage     bool        `json:"hasNextPage"`
	HasPreviousPage bool        `json:"hasPreviousPage"`
}

type connection struct {
	Nodes    []mongoke.Map `json:"nodes"`
	Edges    []edge        `json:"edges"`
	PageInfo pageInfo      `json:"pageInfo"`
}

type edge struct {
	Node   mongoke.Map `json:"node"`
	Cursor interface{} `json:"cursor"`
}

func QueryTypeNodesField(p CreateFieldParams) (*graphql.Field, error) {
	indexableNames := takeIndexableTypeNames(p.SchemaConfig)
	resolver := func(params graphql.ResolveParams) (interface{}, error) {
		args := params.Args
		pagination := paginationFromArgs(args)
		decodedArgs := typeNodesArgs{
			Direction:   mongoke.DESC,
			CursorField: getDefaultCursorField(p.ReturnType, indexableNames),
			Pagination:  pagination,
		}
		err := mapstructure.Decode(args, &decodedArgs)
		if err != nil {
			return nil, err
		}
		decodedArgs, err = addTypeNodesArgsDefaults(decodedArgs)
		if err != nil {
			return nil, err
		}
		if p.InitialWhere != nil {
			mergo.Merge(&decodedArgs.Where, p.InitialWhere)
		}
		opts, err := createFindManyParamsFromArgs(decodedArgs, p.Collection, p.Config.DatabaseUri)
		if err != nil {
			return nil, err
		}
		nodes, err := p.Config.DatabaseFunctions.FindMany(
			opts,
		)
		if err != nil {
			return nil, err
		}

		if len(p.Permissions) == 0 {
			connection := makeConnection(
				nodes,
				decodedArgs.Pagination,
				decodedArgs.CursorField,
			)
			return connection, nil
		}

		jwt := getJwt(params)
		var accessibleNodes []mongoke.Map
		for _, document := range nodes {
			node, err := applyGuardsOnDocument(applyGuardsOnDocumentParams{
				document:  document,
				guards:    p.Permissions,
				jwt:       jwt,
				operation: mongoke.Operations.READ,
			})
			if err != nil {
				// println("got an error while calling applyGuardsOnDocument on findManyField for " + conf.returnType.PrivateName)
				// fmt.Println(err)
				continue
			}
			if node != nil {
				accessibleNodes = append(accessibleNodes, node.(mongoke.Map))
			}
		}
		connection := makeConnection(
			accessibleNodes,
			decodedArgs.Pagination,
			decodedArgs.CursorField,
		)
		// document, err := mongoke.database.findOne()
		// testutil.PrettyPrint(args)
		return connection, nil
	}

	whereArg, err := types.GetWhereArg(p.Config.Cache, indexableNames, p.ReturnType)
	if err != nil {
		return nil, err
	}
	connectionType, err := types.GetConnectionType(p.Config.Cache, p.ReturnType)
	if err != nil {
		return nil, err
	}
	indexableFieldsEnum, err := types.GetIndexableFieldsEnum(p.Config.Cache, indexableNames, p.ReturnType)
	if err != nil {
		return nil, err
	}
	args := graphql.FieldConfigArgument{
		"first":       &graphql.ArgumentConfig{Type: graphql.Int},
		"last":        &graphql.ArgumentConfig{Type: graphql.Int},
		"after":       &graphql.ArgumentConfig{Type: types.AnyScalar},
		"before":      &graphql.ArgumentConfig{Type: types.AnyScalar},
		"direction":   &graphql.ArgumentConfig{Type: types.DirectionEnum},
		"cursorField": &graphql.ArgumentConfig{Type: indexableFieldsEnum},
	}
	if !p.OmitWhere {
		args["where"] = &graphql.ArgumentConfig{Type: whereArg}
	}
	field := graphql.Field{
		Type:    connectionType,
		Args:    args,
		Resolve: resolver,
	}
	return &field, nil
}

func paginationFromArgs(args interface{}) mongoke.Pagination {
	var pag mongoke.Pagination
	err := mapstructure.Decode(args, &pag)
	if err != nil {
		fmt.Println(err)
		return mongoke.Pagination{}
	}
	// increment nodes count so createConnection knows how to set `hasNextPage`
	if pag.First != 0 {
		pag.First++
	}
	if pag.Last != 0 {
		pag.Last++
	}
	// testutil.PrettyPrint(pag)
	return pag
}

func makeConnection(nodes []mongoke.Map, pagination mongoke.Pagination, cursorField string) connection {
	if len(nodes) == 0 {
		return connection{}
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
		nodes = mongoke.ReverseMaps(nodes)
		hasPrev = len(nodes) == int(pagination.Last)
		if hasPrev {
			nodes = nodes[1:]
		}
	}
	if len(nodes) != 0 {
		endCursor = nodes[len(nodes)-1][cursorField]
		startCursor = nodes[0][cursorField]
	}
	return connection{
		Nodes: nodes,
		Edges: makeEdges(nodes, cursorField),
		PageInfo: pageInfo{
			StartCursor:     startCursor,
			EndCursor:       endCursor,
			HasNextPage:     hasNext,
			HasPreviousPage: hasPrev,
		},
	}
}

func makeEdges(nodes []mongoke.Map, cursorField string) []edge {
	edges := make([]edge, len(nodes))
	for _, node := range nodes {
		edges = append(edges, edge{
			Node:   node,
			Cursor: node[cursorField],
		})
	}
	return edges
}

func getJwt(params graphql.ResolveParams) jwt.MapClaims {
	root := params.Info.RootValue
	rootMap, ok := root.(mongoke.Map)
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

const (
	DEFAULT_NODES_COUNT = 20
	MAX_NODES_COUNT     = 40
)

func addTypeNodesArgsDefaults(p typeNodesArgs) (typeNodesArgs, error) {
	if p.Direction == 0 {
		p.Direction = mongoke.ASC
	}
	if p.CursorField == "" {
		p.CursorField = "_id"
	}
	after := p.Pagination.After
	before := p.Pagination.Before
	last := p.Pagination.Last
	first := p.Pagination.First

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
	if after != "" && first == 0 && before == "" {
		return p, errors.New("need `first` or `before` if using `after`")
	}
	if before != "" && (last == 0 && after == "") {
		return p, errors.New("need `last` or `after` if using `before`")
	}
	if first != 0 && last != 0 {
		return p, errors.New("cannot use `first` and `last` together")
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
	return p, nil
}

func createFindManyParamsFromArgs(p typeNodesArgs, collection string, databaseUri string) (mongoke.FindManyParams, error) {
	last := p.Pagination.Last
	first := p.Pagination.First

	opts := mongoke.FindManyParams{
		Collection:  collection,
		DatabaseUri: databaseUri,
		Where:       p.Where,
		Limit:       p.Pagination.First,
		OrderBy: map[string]int{
			p.CursorField: p.Direction,
		},
	}

	// sort order
	sorting := p.Direction
	if last != 0 {
		sorting = -p.Direction
	}
	opts.OrderBy = map[string]int{p.CursorField: sorting}

	// limit
	if last != 0 {
		opts.Limit = int(min(MAX_NODES_COUNT, last))
	}
	if first != 0 {
		opts.Limit = int(min(MAX_NODES_COUNT, first))
	}
	if first == 0 && last == 0 { // when using `after` and `before`
		opts.Limit = int(MAX_NODES_COUNT)
	}
	return opts, nil
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

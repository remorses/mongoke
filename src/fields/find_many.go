package fields

import (
	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
	goke "github.com/remorses/goke/src"
	"github.com/remorses/goke/src/types"
)

func FindMany(p CreateFieldParams) (*graphql.Field, error) {
	resolver := func(params graphql.ResolveParams) (interface{}, error) {
		args := params.Args
		opts := goke.FindManyParams{
			Collection: p.Collection,
		}
		err := mapstructure.Decode(args, &opts)
		if err != nil {
			return nil, err
		}
		if args["where"] != nil {
			where, err := goke.MakeWhereTree(args["where"].(map[string]interface{}), p.InitialWhere)
			if err != nil {
				return nil, err
			}
			opts.Where = where
		}
		nodes, err := p.Config.DatabaseFunctions.FindMany(
			params.Context, opts, func(document goke.Map) (goke.Map, error) {
				// TODO implement check
				return document, nil
			},
		)
		if err != nil {
			return nil, err
		}

		if len(p.Permissions) == 0 {
			return nodes, nil
		}

		jwt := getJwt(params)
		var accessibleNodes []goke.Map
		// TODO move validation logic in a function that returns only accessible documents
		for _, document := range nodes {
			node, err := applyGuardsOnDocument(applyGuardsOnDocumentParams{
				document:  document,
				guards:    p.Permissions,
				jwt:       jwt,
				operation: goke.Operations.READ,
			})
			if err != nil {
				continue
			}
			if node != nil {
				accessibleNodes = append(accessibleNodes, node)
			}
		}
		return accessibleNodes, nil
	}
	indexableNames := takeIndexableTypeNames(p.SchemaConfig)
	whereArg, err := types.GetWhereArg(p.Config.Cache, indexableNames, p.ReturnType)
	if err != nil {
		return nil, err
	}
	orderBy, err := types.GetOrderByArg(p.Config.Cache, indexableNames, p.ReturnType)
	if err != nil {
		return nil, err
	}
	args := graphql.FieldConfigArgument{
		"limit":   &graphql.ArgumentConfig{Type: graphql.Int},
		"offset":  &graphql.ArgumentConfig{Type: graphql.Int},
		"orderBy": &graphql.ArgumentConfig{Type: orderBy},
	}
	if !p.OmitWhere {
		args["where"] = &graphql.ArgumentConfig{Type: whereArg}
	}
	field := graphql.Field{
		Type:    graphql.NewList(p.ReturnType),
		Args:    args,
		Resolve: resolver,
	}
	return &field, nil
}

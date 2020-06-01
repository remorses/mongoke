package fields

import (
	"github.com/graphql-go/graphql"
	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
	mongoke "github.com/remorses/mongoke/src"
	"github.com/remorses/mongoke/src/types"
)

func QueryTypeListField(p CreateFieldParams) (*graphql.Field, error) {
	resolver := func(params graphql.ResolveParams) (interface{}, error) {
		args := params.Args
		opts := mongoke.FindManyParams{
			Collection: p.Collection,
		}
		err := mapstructure.Decode(args, &opts)
		if err != nil {
			return nil, err
		}
		if p.InitialWhere != nil {
			err = mergo.Merge(&opts, p.InitialWhere)
			if err != nil {
				return nil, err
			}
		}
		nodes, err := p.Config.DatabaseFunctions.FindMany(
			params.Context, opts,
		)
		if err != nil {
			return nil, err
		}

		if len(p.Permissions) == 0 {
			return nodes, nil
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
				continue
			}
			if node != nil {
				accessibleNodes = append(accessibleNodes, node.(mongoke.Map))
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

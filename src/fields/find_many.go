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
		opts.Where, err = makeWhere(args, p.InitialWhere, params.Source)
		if err != nil {
			return nil, err
		}
		nodes, err := p.Config.DatabaseFunctions.FindMany(
			params.Context, opts, func(document goke.Map) (goke.Map, error) {
				// println(testutil.Pretty((params.Info.RootValue)))
				return applyGuardsOnDocument(applyGuardsOnDocumentParams{
					isAdmin:            getIsAdmin(params),
					jwt:                getJwt(params),
					defaultPermissions: p.Config.DefaultPermissions,
					document:           document,
					guards:             p.Permissions,
					operation:          goke.Operations.READ,
				})
			},
		)
		if err != nil {
			return nil, err
		}
		return nodes, nil
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

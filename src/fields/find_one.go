package fields

import (
	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
	goke "github.com/remorses/goke/src"
	"github.com/remorses/goke/src/types"
)

func FindOne(p CreateFieldParams) (*graphql.Field, error) {
	indexableNames := takeIndexableTypeNames(p.SchemaConfig)
	resolver := func(params graphql.ResolveParams) (interface{}, error) {
		args := params.Args
		opts := goke.FindManyParams{
			Collection: p.Collection,
			OrderBy: map[string]int{
				getDefaultCursorField(p.ReturnType, indexableNames): goke.DESC,
			},
			Limit: 1,
		}
		err := mapstructure.Decode(args, &opts)
		if err != nil {
			return nil, err
		}
		opts.Where, err = makeWhere(args, p.InitialWhere, params.Source)
		if err != nil {
			return nil, err
		}
		documents, err := p.Config.DatabaseFunctions.FindMany(
			params.Context,
			opts,
			func(document goke.Map) (goke.Map, error) {
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
		if len(documents) == 0 {
			return nil, nil
		}
		return documents[0], nil
	}

	whereArg, err := types.GetWhereArg(p.Config.Cache, indexableNames, p.ReturnType)
	if err != nil {
		return nil, err
	}
	args := graphql.FieldConfigArgument{}
	if !p.OmitWhere {
		args["where"] = &graphql.ArgumentConfig{Type: whereArg}
	}
	field := graphql.Field{
		Type:    p.ReturnType,
		Args:    args,
		Resolve: resolver,
	}
	return &field, nil
}

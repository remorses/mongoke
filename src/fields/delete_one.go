package fields

import (
	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
	goke "github.com/remorses/goke/src"
	"github.com/remorses/goke/src/types"
)

func DeleteOne(p CreateFieldParams) (*graphql.Field, error) {
	resolver := func(params graphql.ResolveParams) (interface{}, error) {
		args := params.Args
		opts := goke.DeleteManyParams{
			Collection: p.Collection,
			Limit:      1,
		}
		err := mapstructure.Decode(args, &opts)
		if err != nil {
			return goke.NodeMutationPayload{}, err
		}
		if args["where"] != nil {
			where, err := goke.MakeWhereTree(args["where"].(map[string]interface{}))
			if err != nil {
				return goke.NodeMutationPayload{}, err
			}
			opts.Where = where
		}
		res, err := p.Config.DatabaseFunctions.DeleteMany(
			params.Context,
			opts,
			func(document goke.Map) (goke.Map, error) {
				return applyGuardsOnDocument(applyGuardsOnDocumentParams{
					isAdmin:            getIsAdmin(params),
					jwt:                getJwt(params),
					defaultPermissions: p.Config.DefaultPermissions,
					document:           document,
					guards:             p.Permissions,
					operation:          goke.Operations.DELETE,
				})
			},
		)
		if err != nil {
			return goke.NodeMutationPayload{}, err
		}
		if len(res.Returning) == 0 {
			return goke.NodeMutationPayload{
				AffectedCount: res.AffectedCount,
			}, nil
		}
		return goke.NodeMutationPayload{
			AffectedCount: res.AffectedCount,
			Returning:     res.Returning[0],
		}, nil
	}
	indexableNames := takeIndexableTypeNames(p.SchemaConfig)
	whereArg, err := types.GetWhereArg(p.Config.Cache, indexableNames, p.ReturnType)
	if err != nil {
		return nil, err
	}
	args := graphql.FieldConfigArgument{
		// "limit":   &graphql.ArgumentConfig{Type: graphql.Int},
		"where": &graphql.ArgumentConfig{Type: whereArg},
	}
	returnType, err := types.GetMutationNodePayload(p.Config.Cache, p.ReturnType)
	if err != nil {
		return nil, err
	}
	field := graphql.Field{
		Type:    returnType,
		Args:    args,
		Resolve: resolver,
	}
	return &field, nil
}

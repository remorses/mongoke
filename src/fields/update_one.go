package fields

import (
	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
	goke "github.com/remorses/goke/src"
	"github.com/remorses/goke/src/types"
)

/*

insertUser(data: {name: "xxx"}) {
	_id
	name
}

*/

func UpdateOne(p CreateFieldParams) (*graphql.Field, error) {
	indexableNames := takeIndexableTypeNames(p.SchemaConfig)
	resolver := func(params graphql.ResolveParams) (interface{}, error) {
		args := params.Args
		opts := goke.UpdateParams{
			Collection: p.Collection,
			Limit:      1,
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
		res, err := p.Config.DatabaseFunctions.UpdateMany(
			params.Context,
			opts,
			func(document goke.Map) (goke.Map, error) {
				return applyGuardsOnDocument(applyGuardsOnDocumentParams{
					jwt:                getJwt(params),
					defaultPermissions: p.Config.DefaultPermissions,
					document:           document,
					guards:             p.Permissions,
					operation:          goke.Operations.UPDATE,
				})
			},
		)
		// println(testutil.Pretty(res))
		if err != nil {
			return nil, err
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

	// if err != nil {
	// 	return nil, err
	// }
	args := graphql.FieldConfigArgument{}
	args["set"] = &graphql.ArgumentConfig{
		Type: types.GetSetUpdateArgument(p.Config.Cache, p.ReturnType),
	}
	whereArg, err := types.GetWhereArg(p.Config.Cache, indexableNames, p.ReturnType)
	if err != nil {
		return nil, err
	}
	args["where"] = &graphql.ArgumentConfig{Type: whereArg}
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

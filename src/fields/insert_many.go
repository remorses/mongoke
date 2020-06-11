package fields

import (
	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
	goke "github.com/remorses/goke/src"
	"github.com/remorses/goke/src/types"
)

/*

insertOneUser(data: {name: "xxx"}) {
	_id
	name
}

*/

func InsertMany(p CreateFieldParams) (*graphql.Field, error) {
	resolver := func(params graphql.ResolveParams) (interface{}, error) {
		args := params.Args
		opts := goke.InsertManyParams{
			Collection: p.Collection,
		}
		err := mapstructure.Decode(args, &opts)
		if err != nil {
			return nil, err
		}
		// TODO insert only nodes the user can insert, based on expressions
		res, err := p.Config.DatabaseFunctions.InsertMany(
			params.Context,
			opts,
			func(document goke.Map) (goke.Map, error) {
				// here the query errors when one document cannot be inserted
				return applyGuardsOnDocument(applyGuardsOnDocumentParams{
					jwt:       getJwt(params),
					document:  document,
					guards:    p.Permissions,
					operation: goke.Operations.CREATE,
				})
			},
		)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	// if err != nil {
	// 	return nil, err
	// }
	args := graphql.FieldConfigArgument{}
	args["data"] = &graphql.ArgumentConfig{
		Type: graphql.NewNonNull(
			graphql.NewList(graphql.NewNonNull(types.TransformToInput(p.Config.Cache, p.ReturnType))),
		),
	}
	returnType, err := types.GetMutationNodesPayload(p.Config.Cache, p.ReturnType)
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

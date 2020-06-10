package fields

import (
	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
	mongoke "github.com/remorses/mongoke/src"
	"github.com/remorses/mongoke/src/testutil"
	"github.com/remorses/mongoke/src/types"
)

/*

insertUser(data: {name: "xxx"}) {
	_id
	name
}

*/

func MutationUpdateNodes(p CreateFieldParams) (*graphql.Field, error) {
	indexableNames := takeIndexableTypeNames(p.SchemaConfig)
	resolver := func(params graphql.ResolveParams) (interface{}, error) {
		args := params.Args
		opts := mongoke.UpdateParams{
			Collection: p.Collection,
		}
		err := mapstructure.Decode(args, &opts)
		if err != nil {
			return nil, err
		}
		if args["where"] != nil {
			where, err := mongoke.MakeWhereTree(args["where"].(map[string]interface{}), p.InitialWhere)
			if err != nil {
				return nil, err
			}
			opts.Where = where
		}
		// TODO update only nodes the user can insert, based on expressions
		res, err := p.Config.DatabaseFunctions.UpdateMany(
			params.Context, opts,
		)
		println(testutil.Pretty(res))
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	// if err != nil {
	// 	return nil, err
	// }
	args := graphql.FieldConfigArgument{}
	args["set"] = &graphql.ArgumentConfig{
		Type: types.GetSetUpdateArgument(p.Config.Cache, p.ReturnType),
	}
	// TODO add update increment field
	// TODO add update append field
	// TODO add update prepend field
	whereArg, err := types.GetWhereArg(p.Config.Cache, indexableNames, p.ReturnType)
	if err != nil {
		return nil, err
	}
	args["where"] = &graphql.ArgumentConfig{Type: whereArg}
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

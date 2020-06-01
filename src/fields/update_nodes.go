package fields

import (
	"github.com/graphql-go/graphql"
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
		// TODO implement resolver
		return nil, nil
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

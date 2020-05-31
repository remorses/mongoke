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

func MutationUpdateOne(p CreateFieldParams) (*graphql.Field, error) {
	indexableNames := takeIndexableTypeNames(p.SchemaConfig)
	resolver := func(params graphql.ResolveParams) (interface{}, error) {
		// TODO implement resolver
		return nil, nil
	}

	// if err != nil {
	// 	return nil, err
	// }
	args := graphql.FieldConfigArgument{}
	args["data"] = &graphql.ArgumentConfig{
		Type: types.TransformToInput(p.Config.Cache, p.ReturnType), // TODO make all partial fields
	}
	whereArg, err := types.GetWhereArg(p.Config.Cache, indexableNames, p.ReturnType)
	if err != nil {
		return nil, err
	}
	args["where"] = &graphql.ArgumentConfig{Type: whereArg}
	args["upsert"] = &graphql.ArgumentConfig{Type: graphql.Boolean}
	field := graphql.Field{
		Type:    p.ReturnType,
		Args:    args,
		Resolve: resolver,
	}
	return &field, nil
}

package fields

import (
	"github.com/graphql-go/graphql"
	"github.com/remorses/mongoke/src/types"
)

/*

insertOneUser(data: {name: "xxx"}) {
	_id
	name
}

*/

func MutationInsertNodes(p CreateFieldParams) (*graphql.Field, error) {
	resolver := func(params graphql.ResolveParams) (interface{}, error) {
		// TODO implement resolver
		return nil, nil
	}

	// if err != nil {
	// 	return nil, err
	// }
	args := graphql.FieldConfigArgument{}
	args["data"] = &graphql.ArgumentConfig{
		Type: graphql.NewNonNull(
			graphql.NewList(types.TransformToInput(p.Config.Cache, p.ReturnType)),
		),
	}
	field := graphql.Field{
		Type:    graphql.NewList(p.ReturnType),
		Args:    args,
		Resolve: resolver,
	}
	return &field, nil
}
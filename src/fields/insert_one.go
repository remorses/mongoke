package fields

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
	mongoke "github.com/remorses/mongoke/src"
	"github.com/remorses/mongoke/src/types"
)

/*

insertUser(data: {name: "xxx"}) {
	_id
	name
}

*/

func MutationInsertOne(p CreateFieldParams) (*graphql.Field, error) {
	resolver := func(params graphql.ResolveParams) (interface{}, error) {
		args := params.Args
		data := mongoke.Map{}
		err := mapstructure.Decode(args["data"], &data)
		if err != nil {
			return nil, err
		}
		if data == nil {
			return nil, errors.New("cannot insert null object")
		}
		opts := mongoke.InsertManyParams{
			Collection: p.Collection,
			Data:       []mongoke.Map{data},
		}

		// TODO insert only nodes the user can insert, based on expressions
		nodes, err := p.Config.DatabaseFunctions.InsertMany(
			params.Context, opts,
		)
		if err != nil {
			return nil, err
		}
		return mongoke.Map{"returning": nodes[0], "affectedCount": 1}, nil
	}

	// if err != nil {
	// 	return nil, err
	// }
	args := graphql.FieldConfigArgument{}
	args["data"] = &graphql.ArgumentConfig{
		Type: graphql.NewNonNull(types.TransformToInput(p.Config.Cache, p.ReturnType)),
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

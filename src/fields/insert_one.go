package fields

import (
	"errors"

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

func InsertOne(p CreateFieldParams) (*graphql.Field, error) {
	resolver := func(params graphql.ResolveParams) (interface{}, error) {
		args := params.Args
		data := goke.Map{}
		err := mapstructure.Decode(args["data"], &data)
		if err != nil {
			return nil, err
		}
		if data == nil {
			return nil, errors.New("cannot insert null object")
		}
		opts := goke.InsertManyParams{
			Collection: p.Collection,
			Data:       []goke.Map{data},
		}

		// TODO insert only nodes the user can insert, based on expressions
		res, err := p.Config.DatabaseFunctions.InsertMany(
			params.Context, opts,
		)
		if err != nil {
			return nil, err
		}
		if len(res.Returning) == 0 {
			return nil, errors.New("could not insert document")
		}
		return goke.Map{"returning": res.Returning[0], "affectedCount": res.AffectedCount}, nil
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

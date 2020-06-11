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
		if args["where"] != nil {
			where, err := goke.MakeWhereTree(args["where"].(map[string]interface{}), p.InitialWhere)
			if err != nil {
				return nil, err
			}
			opts.Where = where
		}
		documents, err := p.Config.DatabaseFunctions.FindMany(params.Context, opts)
		if err != nil {
			return nil, err
		}
		jwt := getJwt(params)
		// don't compute permissions if document is nil
		if len(documents) == 0 {
			return nil, nil
		}
		document := documents[0]
		result, err := applyGuardsOnDocument(applyGuardsOnDocumentParams{
			document:  document,
			guards:    p.Permissions,
			jwt:       jwt,
			operation: goke.Operations.READ,
		})
		if err != nil {
			return nil, err
		}
		return result, nil
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

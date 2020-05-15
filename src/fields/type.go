package fields

import (
	"github.com/graphql-go/graphql"
	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
	mongoke "github.com/remorses/mongoke/src"
	"github.com/remorses/mongoke/src/types"
)

func QueryTypeField(p CreateFieldParams) (*graphql.Field, error) {
	resolver := func(params graphql.ResolveParams) (interface{}, error) {
		args := params.Args
		opts := mongoke.FindManyParams{
			Collection:  p.Collection,
			DatabaseUri: p.Config.DatabaseUri,
			OrderBy:     map[string]int{"_id": mongoke.DESC}, // TODO change _id default based on doc field
			Limit:       1,
		}
		err := mapstructure.Decode(args, &opts)
		if err != nil {
			return nil, err
		}
		if p.InitialWhere != nil {
			mergo.Merge(&opts.Where, p.InitialWhere)
		}
		documents, err := p.Config.DatabaseFunctions.FindMany(opts)
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
			operation: mongoke.Operations.READ,
		})
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	indexableNames := takeIndexableTypeNames(p.SchemaConfig)
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

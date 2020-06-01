package schema

import (
	"github.com/graphql-go/graphql"
	mongoke "github.com/remorses/mongoke/src"
	"github.com/remorses/mongoke/src/fields"
)

func makeMutation(Config mongoke.Config, baseSchemaConfig graphql.SchemaConfig) (*graphql.Object, error) {
	mutationFields := graphql.Fields{}
	for _, gqlType := range baseSchemaConfig.Types {
		var object graphql.Type
		switch v := gqlType.(type) {
		case *graphql.Object, *graphql.Union, *graphql.Interface:
			object = v
		default:
			continue
		}
		typeConf := Config.GetTypeConfig(gqlType.Name())

		if unexposedType(typeConf) {
			println("ignoring not exposed type " + gqlType.Name())
			continue
		}
		p := fields.CreateFieldParams{
			Config:       Config,
			ReturnType:   object,
			Permissions:  typeConf.Permissions,
			Collection:   typeConf.Collection,
			SchemaConfig: baseSchemaConfig,
		}
		// insertOne
		insertOne, err := fields.MutationInsertOne(p)
		if err != nil {
			return nil, err
		}
		mutationFields["insert"+object.Name()] = insertOne
		// insertMany
		insertNodes, err := fields.MutationInsertNodes(p)
		if err != nil {
			return nil, err
		}
		mutationFields["insert"+object.Name()+"Nodes"] = insertNodes
		// updateOne
		updateOne, err := fields.MutationUpdateOne(p)
		if err != nil {
			return nil, err
		}
		mutationFields["update"+object.Name()] = updateOne
		updateNodes, err := fields.MutationUpdateNodes(p)
		if err != nil {
			return nil, err
		}
		mutationFields["update"+object.Name()+"Nodes"] = updateNodes
	}
	mutation := graphql.NewObject(graphql.ObjectConfig{Name: "Mutation", Fields: mutationFields})
	return mutation, nil
}

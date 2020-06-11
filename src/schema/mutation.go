package schema

import (
	"github.com/graphql-go/graphql"
	goke "github.com/remorses/goke/src"
	"github.com/remorses/goke/src/fields"
)

func makeMutation(Config goke.Config, baseSchemaConfig graphql.SchemaConfig) (*graphql.Object, error) {
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
			// println("ignoring not exposed type " + gqlType.Name())
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
		insertOne, err := fields.InsertOne(p)
		if err != nil {
			return nil, err
		}
		mutationFields["insertOne"+object.Name()] = insertOne

		// insertMany
		insertMany, err := fields.InsertMany(p)
		if err != nil {
			return nil, err
		}
		mutationFields["insertMany"+object.Name()] = insertMany

		// updateOne
		updateOne, err := fields.UpdateOne(p)
		if err != nil {
			return nil, err
		}
		mutationFields["updateOne"+object.Name()] = updateOne

		// updateMany
		updateMany, err := fields.UpdateMany(p)
		if err != nil {
			return nil, err
		}
		mutationFields["updateMany"+object.Name()] = updateMany

		// deleteOne
		deleteOne, err := fields.DeleteMany(p)
		if err != nil {
			return nil, err
		}
		mutationFields["deleteOne"+object.Name()] = deleteOne

		// deleteMany
		deleteMany, err := fields.DeleteMany(p)
		if err != nil {
			return nil, err
		}
		mutationFields["deleteMany"+object.Name()] = deleteMany
	}
	mutation := graphql.NewObject(graphql.ObjectConfig{Name: "Mutation", Fields: mutationFields})
	return mutation, nil
}

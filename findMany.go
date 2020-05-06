package mongoke

import (
	"github.com/graphql-go/graphql"
)

type findManyResolverConfig struct {
	collection string
	database   string
	returnType *graphql.Object
}

var directionEnum = graphql.NewEnum(graphql.EnumConfig{
	Name:        "Direction",
	Description: "asc or desc",
	Values: graphql.EnumValueConfigMap{
		"ASC": &graphql.EnumValueConfig{
			Value:       1,
			Description: "ascending",
		},
		"DESC": &graphql.EnumValueConfig{
			Value:       2,
			Description: "Descending",
		},
	},
})

func findManyResolver(conf findManyResolverConfig) *graphql.Field {
	// TODO create the where argument based on the object fields

	resolver := func(params graphql.ResolveParams) (interface{}, error) {
		args := params.Args
		// TODO get item from database
		// check authorization guards
		// if interface or union set the right __typeName
		prettyPrint(args)
		return "world", nil
	}
	return &graphql.Field{
		Type: graphql.NewObject(connectionType(conf.returnType)),
		Args: graphql.FieldConfigArgument{
			"where":     whereArgument(*conf.returnType),
			"first":     &graphql.ArgumentConfig{Type: graphql.Int},
			"last":      &graphql.ArgumentConfig{Type: graphql.Int},
			"direction": &graphql.ArgumentConfig{Type: directionEnum},
		},
		Resolve: resolver,
	}
}

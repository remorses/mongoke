package mongoke

import (
	"github.com/graphql-go/graphql"
)

type findOneFieldConfig struct {
	collection string
	database   string
	returnType *graphql.Object
}

func (mongoke *Mongoke) findOneField(conf findOneFieldConfig) *graphql.Field {
	// TODO create the where argument based on the object fields

	resolver := func(params graphql.ResolveParams) (interface{}, error) {
		args := params.Args
		// TODO get item from database
		// check authorization guards
		// if interface or union set the right __typeName
		prettyPrint(args)
		return "world", nil
	}
	whereArg, err := mongoke.getWhereArg(conf.returnType)
	if err != nil {
		panic(err)
	}
	return &graphql.Field{
		Type: conf.returnType,
		Args: graphql.FieldConfigArgument{
			"where": &graphql.ArgumentConfig{Type: whereArg},
		},
		Resolve: resolver,
	}
}

type findManyFieldConfig struct {
	collection string
	database   string
	returnType *graphql.Object
}

func (mongoke *Mongoke) findManyField(conf findManyFieldConfig) *graphql.Field {
	resolver := func(params graphql.ResolveParams) (interface{}, error) {
		args := params.Args
		// TODO get item from database
		// check authorization guards
		// if interface or union set the right __typeName
		prettyPrint(args)
		return "world", nil
	}
	whereArg, err := mongoke.getWhereArg(conf.returnType)
	if err != nil {
		panic(err)
	}
	connectionType, err := mongoke.getConnectionType(conf.returnType)
	if err != nil {
		panic(err)
	}
	return &graphql.Field{
		Type: connectionType,
		Args: graphql.FieldConfigArgument{
			"where":     &graphql.ArgumentConfig{Type: whereArg},
			"first":     &graphql.ArgumentConfig{Type: graphql.Int},
			"last":      &graphql.ArgumentConfig{Type: graphql.Int},
			"direction": &graphql.ArgumentConfig{Type: directionEnum},
		},
		Resolve: resolver,
	}
}

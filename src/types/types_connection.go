package types

import (
	"github.com/graphql-go/graphql"
	mongoke "github.com/remorses/mongoke/src"
)

func MakeConnectionTypeName(object graphql.Type) string {
	return object.Name() + "Connection"
}

func GetConnectionType(cache mongoke.Map, object graphql.Type) (*graphql.Object, error) {
	name := MakeConnectionTypeName(object)
	// get cached value to not dupe
	if item, ok := cache[name].(*graphql.Object); ok {
		return item, nil
	}
	edgeNode := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        object.Name() + "Edge",
			Description: "Edge",
			Fields: graphql.Fields{
				"node": &graphql.Field{
					Type: object,
				},
				"cursor": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)
	fields := graphql.Fields{
		"nodes": &graphql.Field{
			Type: graphql.NewList(object),
		},
		"edges": &graphql.Field{
			Type: graphql.NewList(edgeNode),
		},
		"pageInfo": &graphql.Field{
			Type: PageInfo,
		},
	}
	connection := graphql.NewObject(graphql.ObjectConfig{
		Name:   name,
		Fields: fields,
	})
	cache[name] = connection
	return connection, nil
}

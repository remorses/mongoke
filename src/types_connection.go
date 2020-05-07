package mongoke

import (
	"github.com/graphql-go/graphql"
)

func connectionType(object *graphql.Object) graphql.ObjectConfig {
	name := object.Name() + "Connection"
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
			Type: pageInfo,
		},
	}
	connection := graphql.ObjectConfig{
		Name:   name,
		Fields: fields,
	}
	return connection
}

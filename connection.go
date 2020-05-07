package mongoke

import (
	"github.com/graphql-go/graphql"
)

var pageInfo = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "PageInfo",
		Description: "Pagination information",
		Fields: graphql.Fields{
			"startCursor": &graphql.Field{
				Type: graphql.String, // TODO should be anyscalar
			},
			"endCursor": &graphql.Field{
				Type: graphql.String,
			},
			"hasNextPage": &graphql.Field{
				Type: graphql.Boolean,
			},
			"hasPreviousPage": &graphql.Field{
				Type: graphql.Boolean,
			},
		},
	},
)

func connectionType(object *graphql.Object) graphql.ObjectConfig {
	name := object.Name() + "Connection"
	node := graphql.NewObject(
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
			Type: graphql.NewList(node),
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

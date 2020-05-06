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

func connectionType(object graphql.Object) *graphql.ObjectConfig {
	fields := graphql.Fields{
		"nodes": &graphql.Field{
			Type: graphql.NewList(&object),
		},
		"pageInfo": &graphql.Field{
			Type: pageInfo,
		},
	}
	connection := &graphql.ObjectConfig{
		Name:   object.Name() + "Connection",
		Fields: fields,
	}

	return connection
}

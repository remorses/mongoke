package mongoke

import (
	"github.com/graphql-go/graphql"
)

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

package mongoke

import (
	"github.com/graphql-go/graphql"
)

func whereArgument(fields []*graphql.FieldDefinition) *graphql.ArgumentConfig {
	where := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "Where", Fields: graphql.InputObjectConfigFieldMap{
			"eq": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
		},
	})
	arg := graphql.ArgumentConfig{Type: where}
	return &arg
}

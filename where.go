package mongoke

import (
	"github.com/graphql-go/graphql"
)

func whereArgument(object graphql.Object) *graphql.ArgumentConfig {
	scalars := takeScalarFields(object, []string{}) // TODO add here the enums and scalars names
	inputFields := graphql.InputObjectConfigFieldMap{}
	for _, field := range scalars {
		fieldWhere := fieldWhereArgument(field)
		inputFields[field.Name] = &graphql.InputObjectFieldConfig{
			Type:        fieldWhere,
			Description: "The Mongodb match object for the field " + field.Name,
		}
	}
	where := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:   object.Name() + "Where",
		Fields: inputFields,
	})
	inputFields["or"] = &graphql.InputObjectFieldConfig{Type: where}
	inputFields["and"] = &graphql.InputObjectFieldConfig{Type: where}
	arg := graphql.ArgumentConfig{
		Type:        where,
		Description: "Where type for " + object.Name(),
	}
	return &arg
}

func fieldWhereArgument(field *graphql.FieldDefinition) *graphql.InputObject {
	currentType := &graphql.InputObjectFieldConfig{Type: field.Type}
	fieldWhere := graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name: field.Name + "Where",
			Fields: graphql.InputObjectConfigFieldMap{
				"eq":  currentType,
				"neq": currentType,
				"in": &graphql.InputObjectFieldConfig{
					Type: graphql.NewList(field.Type),
				},
				"nin": &graphql.InputObjectFieldConfig{
					Type: graphql.NewList(field.Type),
				},
			},
		},
	)
	return fieldWhere
}

package mongoke

import (
	"errors"

	"github.com/graphql-go/graphql"
)

func makeWhereArgumentName(object *graphql.Object) string {
	return object.Name() + "Where"
}

func (mongoke *Mongoke) getWhereArg(object *graphql.Object) (*graphql.InputObject, error) {
	name := makeWhereArgumentName(object)
	if item, ok := mongoke.typeMap[name]; ok {
		if t, ok := item.(*graphql.InputObject); ok {
			return t, nil
		}
		return nil, errors.New("cannot cast where type for " + name)
	}
	scalars := takeScalarFields(object, []string{}) // TODO add here the enums and scalars names
	inputFields := graphql.InputObjectConfigFieldMap{}
	for _, field := range scalars {
		fieldWhere := fieldWhereArgument(field, object.PrivateName)
		inputFields[field.Name] = &graphql.InputObjectFieldConfig{
			Type:        fieldWhere,
			Description: "The Mongodb match object for the field " + field.Name,
		}
	}
	where := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:   name,
		Fields: inputFields,
	})
	inputFields["or"] = &graphql.InputObjectFieldConfig{Type: where}
	inputFields["and"] = &graphql.InputObjectFieldConfig{Type: where}
	mongoke.typeMap[name] = where
	return where, nil
}

func fieldWhereArgument(field *graphql.FieldDefinition, parentName string) *graphql.InputObject {
	name := parentName + field.Name + "FieldMatch"
	currentType := &graphql.InputObjectFieldConfig{Type: field.Type}
	fieldWhere := graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name: name,
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

package mongoke

import (
	"github.com/graphql-go/graphql"
)

func makeWhereArgumentName(object graphql.Type) string {
	return object.Name() + "Where"
}

func getWhereArg(cache Map, indexableNames []string, object graphql.Type) (*graphql.InputObject, error) {
	name := makeWhereArgumentName(object)
	if item, ok := cache[name].(*graphql.InputObject); ok {
		return item, nil
	}
	scalars := takeIndexableFields(indexableNames, object)
	inputFields := graphql.InputObjectConfigFieldMap{}
	for _, field := range scalars {
		fieldWhere, err := getFieldWhereArg(cache, field, object.Name())
		if err != nil {
			return nil, err
		}
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
	cache[name] = where
	return where, nil
}

// this is be based on a type, like scalars, enums, ..., cache it in mongoke and replace name
func getFieldWhereArg(cache Map, field *graphql.FieldDefinition, parentName string) (*graphql.InputObject, error) {
	name := field.Type.Name() + "Where"
	if item, ok := cache[name].(*graphql.InputObject); ok {
		return item, nil
	}
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
	cache[name] = fieldWhere
	return fieldWhere, nil
}

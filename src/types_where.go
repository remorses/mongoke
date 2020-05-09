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
	scalars := mongoke.takeIndexableFields(object)
	inputFields := graphql.InputObjectConfigFieldMap{}
	for _, field := range scalars {
		fieldWhere, err := mongoke.getFieldWhereArg(field, object.PrivateName)
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
	mongoke.typeMap[name] = where
	return where, nil
}

// this is be based on a type, like scalars, enums, ..., cache it in mongoke and replace name
func (mongoke *Mongoke) getFieldWhereArg(field *graphql.FieldDefinition, parentName string) (*graphql.InputObject, error) {
	name := field.Type.Name() + "Where"
	if item, ok := mongoke.typeMap[name]; ok {
		if t, ok := item.(*graphql.InputObject); ok {
			return t, nil
		}
		return nil, errors.New("cannot cast field where type for " + name)
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
	mongoke.typeMap[name] = fieldWhere
	return fieldWhere, nil
}

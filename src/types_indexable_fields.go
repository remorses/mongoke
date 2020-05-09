package mongoke

import (
	"errors"

	"github.com/graphql-go/graphql"
)

func makeIndexableFieldsName(object *graphql.Object) string {
	return object.Name() + "IndexableFields"
}

func (mongoke *Mongoke) getIndexableFieldsEnum(object *graphql.Object) (*graphql.Enum, error) {
	name := makeIndexableFieldsName(object)
	if item, ok := mongoke.typeMap[name]; ok {
		if t, ok := item.(*graphql.Enum); ok {
			return t, nil
		}
		return nil, errors.New("cannot cast indexable fields type for " + name)
	}
	scalars := mongoke.takeIndexableFields(object)
	enumValues := graphql.EnumValueConfigMap{}
	for _, field := range scalars {
		enumValues[field.Name] = &graphql.EnumValueConfig{
			Value:       field.Name,
			Description: "The field enum for " + field.Name,
		}
	}
	enum := graphql.NewEnum(graphql.EnumConfig{
		Name:   name,
		Values: enumValues,
	})
	mongoke.typeMap[name] = enum
	return enum, nil
}

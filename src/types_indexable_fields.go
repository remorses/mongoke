package mongoke

import (
	"github.com/graphql-go/graphql"
)

func makeIndexableFieldsName(object graphql.Type) string {
	return object.Name() + "IndexableFields"
}

func getIndexableFieldsEnum(cache Map, indexableNames []string, object graphql.Type) (*graphql.Enum, error) {
	name := makeIndexableFieldsName(object)
	if item, ok := cache[name].(*graphql.Enum); ok {
		return item, nil
	}
	scalars := takeIndexableFields(indexableNames, object)
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
	cache[name] = enum
	return enum, nil
}

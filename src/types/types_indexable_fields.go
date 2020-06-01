package types

import (
	"github.com/graphql-go/graphql"
	mongoke "github.com/remorses/mongoke/src"
)

func GetIndexableFieldsEnum(cache mongoke.Map, indexableNames []string, object graphql.Type) (*graphql.Enum, error) {
	name := object.Name() + "IndexableFields"
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

// indexableFields
func takeIndexableFields(indexableNames []string, object graphql.Type) []*graphql.FieldDefinition {
	indexableFields := make([]*graphql.FieldDefinition, 0)
	for _, v := range GetTypeFields(object) {
		var t graphql.Type
		if n, ok := v.Type.(*graphql.NonNull); ok {
			t = n.OfType
		} else {
			t = v.Type
		}
		// TODO maybe handle list as an indexable type to add the contains operator
		typeName := t.Name()
		f := &graphql.FieldDefinition{Type: t, Name: v.Name, Args: v.Args, Description: v.Description, Resolve: v.Resolve}
		switch typeName {
		case "String", "Boolean", "Int", "Float", "ID", "DateTime":
			indexableFields = append(indexableFields, f)
		}
		if contains(indexableNames, typeName) {
			indexableFields = append(indexableFields, f)
		}
	}
	return indexableFields
}

func GetTypeFields(object graphql.Type) graphql.FieldDefinitionMap {
	fieldMap := graphql.FieldDefinitionMap{}
	switch v := object.(type) {
	case *graphql.Object:
		return v.Fields()
	case *graphql.Interface:
		return v.Fields()
	case *graphql.Union:
		for _, t := range v.Types() {
			for k, field := range t.Fields() {
				fieldMap[k] = field
			}
		}
		return fieldMap
	default:
		return graphql.FieldDefinitionMap{}
	}
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

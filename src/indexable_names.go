package mongoke

import (
	"github.com/graphql-go/graphql"
)

// indexableFields
func takeIndexableFields(indexableNames []string, object graphql.Type) []*graphql.FieldDefinition {
	// indexableNames := takeIndexableTypeNames(schemaConfig)
	indexableFields := make([]*graphql.FieldDefinition, 0)
	for _, v := range getTypeFields(object) {
		typeName := v.Type.Name()
		switch typeName {
		case "String", "Boolean", "Int", "Float", "ID", "DateTime":
			indexableFields = append(indexableFields, v)
		}
		if contains(indexableNames, typeName) {
			indexableFields = append(indexableFields, v)
		}
	}
	return indexableFields
}

func getTypeFields(object graphql.Type) graphql.FieldDefinitionMap {
	fieldMap := graphql.FieldDefinitionMap{}
	switch v := object.(type) {
	case *graphql.Object:
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

// to be used in takeScalarFields
func takeIndexableTypeNames(baseSchemaConfig graphql.SchemaConfig) []string {
	names := make([]string, 0)
	for _, gqlType := range baseSchemaConfig.Types {
		if graphql.IsLeafType(gqlType) {
			names = append(names, gqlType.Name())
		}
	}
	return names
}

func takeEnumTypes(baseSchemaConfig graphql.SchemaConfig) []*graphql.Enum {
	enums := make([]*graphql.Enum, 0)
	for _, gqlType := range baseSchemaConfig.Types {
		enum, ok := gqlType.(*graphql.Enum)
		if !ok {
			continue
		}
		enums = append(enums, enum)
	}
	return enums
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

package mongoke

import (
	"github.com/graphql-go/graphql"
)

// indexableFields
func (m Mongoke) takeIndexableFields(object *graphql.Object) []*graphql.FieldDefinition {
	indexableNames := takeIndexableTypeNames(m.schemaConfig)
	indexableFields := make([]*graphql.FieldDefinition, 0)
	for _, v := range object.Fields() {
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

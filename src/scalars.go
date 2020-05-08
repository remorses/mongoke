package mongoke

import (
	"github.com/graphql-go/graphql"
)

// indexableFields
func takeScalarFields(object *graphql.Object, scalars []string) []*graphql.FieldDefinition {
	scalarFields := make([]*graphql.FieldDefinition, 0)
	for _, v := range object.Fields() {
		typeName := v.Type.Name()
		switch typeName {
		case "String", "Boolean", "Int", "Float", "ID", "DateTime":
			scalarFields = append(scalarFields, v)
		}
		if contains(scalars, typeName) {
			scalarFields = append(scalarFields, v)
		}
	}
	return scalarFields
}

// to be used in takeScalarFields
func takeScalarTypeNames(baseSchemaConfig graphql.SchemaConfig) []string {
	names := make([]string, 0)
	enums := takeEnumTypes(baseSchemaConfig)
	for _, scalar := range append(enums) { // TODO add scalar typeNames to compute indexable fields, i could use graphql.IsLeafType
		names = append(names, scalar.Name())
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

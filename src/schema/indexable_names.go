package schema

import (
	"github.com/graphql-go/graphql"
)

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

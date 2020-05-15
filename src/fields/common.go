package fields

import (
	"github.com/graphql-go/graphql"
	mongoke "github.com/remorses/mongoke/src"
)

type CreateFieldParams struct {
	Config       mongoke.Config
	Collection   string
	InitialWhere map[string]mongoke.Filter
	Permissions  []mongoke.AuthGuard
	ReturnType   graphql.Type
	SchemaConfig graphql.SchemaConfig
	OmitWhere    bool
}

func takeIndexableTypeNames(baseSchemaConfig graphql.SchemaConfig) []string {
	names := make([]string, 0)
	for _, gqlType := range baseSchemaConfig.Types {
		if _, ok := gqlType.(*graphql.List); ok {
			// because isLeafType also returns lists
			continue
		}
		if graphql.IsLeafType(gqlType) {
			names = append(names, gqlType.Name())
		}
	}
	return names
}

func getDefaultCursorField(scalarNames []string) string {
	// TODO customize default cursor field via config
	for _, name := range scalarNames {
		if name == "_id" || name == "id" {
			return name
		}
	}
	if len(scalarNames) != 0 {

		return scalarNames[0]
	}
	return ""
}

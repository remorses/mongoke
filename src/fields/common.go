package fields

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/graphql-go/graphql"
	mongoke "github.com/remorses/mongoke/src"
	"github.com/remorses/mongoke/src/types"
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

func getDefaultCursorField(object graphql.Type, scalarTypes []string) string {
	var scalarNames []string
	fieldNames := types.GetTypeFields(object)
	for name, t := range fieldNames {
		if contains(scalarTypes, t.Type.Name()) {
			scalarNames = append(scalarNames, name)
		}
	}
	for _, name := range scalarNames {
		if name == "_id" || name == "id" { // TODO customize default cursor field via config
			return name
		}
	}
	if len(scalarNames) != 0 {
		return scalarNames[0]
	}
	return ""
}

func getJwt(params graphql.ResolveParams) jwt.MapClaims {
	root := params.Info.RootValue
	rootMap, ok := root.(mongoke.Map)
	if !ok {
		return jwt.MapClaims{}
	}
	v, ok := rootMap["jwt"]
	if !ok {
		return jwt.MapClaims{}
	}
	jwtMap, ok := v.(jwt.MapClaims)
	if !ok {
		return jwt.MapClaims{}
	}
	return jwtMap
}

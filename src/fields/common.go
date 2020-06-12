package fields

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/graphql-go/graphql"
	goke "github.com/remorses/goke/src"
	"github.com/remorses/goke/src/types"
)

type CreateFieldParams struct {
	Config       goke.Config
	Collection   string
	InitialWhere map[string]goke.Filter
	Permissions  []goke.AuthGuard
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
	if root == nil {
		return jwt.MapClaims{}
	}
	// rootMap, ok := root.(goke.Map)
	// if !ok {
	// 	println("not a map")
	// 	return jwt.MapClaims{}
	// }
	rootMap, ok := root.(map[string]interface{})
	if !ok {
		println("WARNING: jwt is not a map")
		return jwt.MapClaims{}
	}
	v, ok := rootMap["jwt"]
	if !ok {
		return jwt.MapClaims{}
	}
	jwtMap, ok := v.(jwt.MapClaims)
	if !ok {
		println("WARNING: jwt is not a MapClaims")
		return jwt.MapClaims{}
	}
	return jwtMap
}

func getIsAdmin(params graphql.ResolveParams) bool {
	root := params.Info.RootValue
	rootMap, ok := root.(goke.Map)
	if !ok {
		return false
	}
	v, ok := rootMap["isAdmin"]
	if !ok {
		return false
	}
	isAdmin, ok := v.(bool)
	if !ok {
		return false
	}
	return isAdmin
}

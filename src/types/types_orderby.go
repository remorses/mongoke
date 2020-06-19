package types

import (
	"github.com/graphql-go/graphql"
	goke "github.com/remorses/goke/src"
)

func MakeOrderByArgumentName(object graphql.Type) string {
	return object.Name() + "OrderBy"
}

func GetOrderByArg(cache goke.Map, indexableNames []string, object graphql.Type) (*graphql.InputObject, error) {
	name := MakeOrderByArgumentName(object)
	if item, ok := cache[name].(*graphql.InputObject); ok {
		return item, nil
	}
	scalars := takeIndexableFields(indexableNames, object)
	inputFields := graphql.InputObjectConfigFieldMap{}
	for _, field := range scalars {
		inputFields[field.Name] = &graphql.InputObjectFieldConfig{
			Type: DirectionEnum,
		}
	}
	orderBy := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:   name,
		Fields: inputFields,
	})
	cache[name] = orderBy
	return orderBy, nil
}

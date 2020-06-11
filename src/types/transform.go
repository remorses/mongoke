package types

import (
	"github.com/graphql-go/graphql"
	goke "github.com/remorses/goke/src"
)

func MakeInputPartial(cache goke.Map, object *graphql.InputObject) *graphql.InputObject {
	name := object.Name() + "Partial"
	if item, ok := cache[name].(*graphql.InputObject); ok {
		return item
	}
	fields := graphql.InputObjectConfigFieldMap{}
	for k, v := range object.Fields() {
		if f, ok := v.Type.(*graphql.NonNull); ok {
			fields[k] = &graphql.InputObjectFieldConfig{
				Type:        f.OfType,
				Description: f.Description(),
			}
		} else {
			fields[k] = &graphql.InputObjectFieldConfig{
				Type:        v.Type,
				Description: v.Description(),
			}
		}
	}
	res := graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name:        name,
			Fields:      fields,
			Description: object.Description(),
		},
	)
	cache[name] = res
	return res
}

func TransformToInput(cache goke.Map, object graphql.Type) *graphql.InputObject {
	name := object.Name() + "Input"
	if item, ok := cache[name].(*graphql.InputObject); ok {
		return item
	}
	fields := graphql.InputObjectConfigFieldMap{}
	for _, field := range GetTypeFields(object) {
		t := field.Type
		if v, ok := t.(*graphql.Object); ok {
			t = TransformToInput(cache, v)
		}
		fields[field.Name] = &graphql.InputObjectFieldConfig{
			Type:        t,
			Description: field.Description,
		}
	}
	config := graphql.InputObjectConfig{
		Name:   name,
		Fields: fields,
	}
	input := graphql.NewInputObject(config)
	cache[name] = input
	return input
}

package types

import (
	"github.com/graphql-go/graphql"
	mongoke "github.com/remorses/mongoke/src"
)

func TransformToInput(cache mongoke.Map, t graphql.Type) graphql.Type {
	switch t.(type) {
	case *graphql.Object:
		return objectToInputObject(cache, t.(*graphql.Object))
	} // TODO transform unions to input objects
	return t
}

func makeObjectInputName(t graphql.Type) string {
	return t.Name() + "Input"
}

func objectToInputObject(cache mongoke.Map, object *graphql.Object) *graphql.InputObject {
	name := makeObjectInputName(object)
	if item, ok := cache[name].(*graphql.InputObject); ok {
		return item
	}
	fields := graphql.InputObjectConfigFieldMap{}
	for _, field := range object.Fields() {
		t := field.Type
		if v, ok := t.(*graphql.Object); ok {
			t = objectToInputObject(cache, v)
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

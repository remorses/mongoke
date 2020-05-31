package types

import "github.com/graphql-go/graphql"

func TransformToInput(t graphql.Type) graphql.Type {
	switch t.(type) {
	case *graphql.Object:
		return objectToInputObject(t.(*graphql.Object))
	} // TODO transform unions to input objects
	return t
}

func objectToInputObject(object *graphql.Object) *graphql.InputObject {
	fields := graphql.InputObjectConfigFieldMap{}
	for _, field := range object.Fields() {
		t := field.Type
		if v, ok := t.(*graphql.Object); ok {
			t = objectToInputObject(v)
		}
		fields[field.Name] = &graphql.InputObjectFieldConfig{
			Type:        t,
			Description: field.Description,
		}
	}
	config := graphql.InputObjectConfig{
		Name:   object.Name() + "Input",
		Fields: fields,
	}
	return graphql.NewInputObject(config)
}

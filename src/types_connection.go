package mongoke

import (
	"errors"

	"github.com/graphql-go/graphql"
)

func makeConnectionTypeName(object *graphql.Object) string {
	return object.Name() + "Connection"
}

func (mongoke *Mongoke) getConnectionType(object *graphql.Object) (*graphql.Object, error) {
	name := makeConnectionTypeName(object)
	// get cached value to not dupe
	if item, ok := mongoke.typeMap[name]; ok {
		if t, ok := item.(*graphql.Object); ok {
			return t, nil
		}
		return nil, errors.New("cannot cast connection type for " + name)
	}
	edgeNode := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        object.Name() + "Edge",
			Description: "Edge",
			Fields: graphql.Fields{
				"node": &graphql.Field{
					Type: object,
				},
				"cursor": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)
	fields := graphql.Fields{
		"nodes": &graphql.Field{
			Type: graphql.NewList(object),
		},
		"edges": &graphql.Field{
			Type: graphql.NewList(edgeNode),
		},
		"pageInfo": &graphql.Field{
			Type: pageInfo,
		},
	}
	connection := graphql.NewObject(graphql.ObjectConfig{
		Name:   name,
		Fields: fields,
	})
	mongoke.typeMap[name] = connection
	return connection, nil
}

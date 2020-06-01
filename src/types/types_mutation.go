package types

import (
	"github.com/graphql-go/graphql"
	mongoke "github.com/remorses/mongoke/src"
)

/*

type MutationPayload {
	affectedDocuments
	returning
}

Mongodb update many returns
type UpdateResult struct {
    MatchedCount  int64       // The number of documents matched by the filter.
    ModifiedCount int64       // The number of documents modified by the operation.
    UpsertedCount int64       // The number of documents upserted by the operation.
    UpsertedID    interface{} // The _id field of the upserted document, or nil if no upsert was done.
}

In firestore i have to handle everything manually, upserts are done manually, i also have to update one by one



in postgres

*/

func GetSetUpdateArgument(cache mongoke.Map, object graphql.Type) *graphql.InputObject {
	name := object.Name() + "SetUpdate"
	if item, ok := cache[name].(*graphql.InputObject); ok {
		return item
	}
	input := MakeInputPartial(cache, TransformToInput(cache, object))
	fields := graphql.InputObjectConfigFieldMap{}
	for k, v := range input.Fields() {
		fields[k] = &graphql.InputObjectFieldConfig{
			Type: v.Type,
		}
	}
	set := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:   name,
		Fields: fields,
	})
	cache[name] = set
	return set
}

func GetMutationNodePayload(cache mongoke.Map, object graphql.Type) (*graphql.Object, error) {
	name := object.Name() + "MutationNodePayload"
	// get cached value to not dupe
	if item, ok := cache[name].(*graphql.Object); ok {
		return item, nil
	}
	fields := graphql.Fields{
		"returning": &graphql.Field{
			Type: object,
		},
		"affectedCount": &graphql.Field{
			Type: graphql.Int,
		},
	}
	payload := graphql.NewObject(graphql.ObjectConfig{
		Name:   name,
		Fields: fields,
	})
	cache[name] = payload
	return payload, nil
}

func GetMutationNodesPayload(cache mongoke.Map, object graphql.Type) (*graphql.Object, error) {
	name := object.Name() + "MutationNodesPayload"
	// get cached value to not dupe
	if item, ok := cache[name].(*graphql.Object); ok {
		return item, nil
	}
	fields := graphql.Fields{
		"returning": &graphql.Field{
			Type: graphql.NewList(object),
		},
		"affectedCount": &graphql.Field{
			Type: graphql.Int,
		},
	}
	payload := graphql.NewObject(graphql.ObjectConfig{
		Name:   name,
		Fields: fields,
	})
	cache[name] = payload
	return payload, nil
}

// func GetUpdateArg(cache mongoke.Map, indexableNames []string, object graphql.Type) (*graphql.InputObject, error) {
// 	name := object.Name() + "Update"
// 	if item, ok := cache[name].(*graphql.InputObject); ok {
// 		return item, nil
// 	}
// 	// scalars := takeIndexableFields(indexableNames, object)
// 	inputFields := graphql.InputObjectConfigFieldMap{}
// 	inputFields["set"] = &graphql.InputObjectFieldConfig{
// 		Type: makeUpdateSetArgument(cache, object),
// 	}

// 	where := graphql.NewInputObject(graphql.InputObjectConfig{
// 		Name:   name,
// 		Fields: inputFields,
// 	})
// 	cache[name] = where
// 	return where, nil
// }

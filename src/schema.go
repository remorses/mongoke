package mongoke

import (
	"errors"

	"github.com/graphql-go/graphql"
)

func (mongoke *Mongoke) generateSchema() (graphql.Schema, error) {
	queryFields := graphql.Fields{}
	mutationFields := graphql.Fields{}

	baseSchemaConfig := mongoke.schemaConfig
	for _, gqlType := range baseSchemaConfig.Types {
		// TODO handle unions types, adding a Field function to Object and creating a shared interface

		switch gqlType.(type) {
		case *graphql.Union:
			return graphql.Schema{}, errors.New("union types are still not supported, for union type" + gqlType.Name())
		// case *graphql.Scalar: // TODO interfaces are not found, throws no definition found for {interfaceName}
		// 	return graphql.Schema{}, errors.New("scalar types are still not supported, for scalar type" + gqlType.Name())
		case *graphql.Interface:
			println("interface")
			return graphql.Schema{}, errors.New("interface types are still not supported, for interface type" + gqlType.Name())
		}

		object, ok := gqlType.(*graphql.Object)
		if !ok {
			continue
		}

		typeConf := mongoke.Config.getTypeConfig(gqlType.Name())

		if typeConf == nil || (typeConf.Exposed != nil && !*typeConf.Exposed) {
			println("ignoring not exposed type " + gqlType.Name())
			continue
		}

		if typeConf.Collection == "" {
			return graphql.Schema{}, errors.New("no collection given for type " + gqlType.Name())
		}

		queryFields[object.Name()] = mongoke.findOneField(
			findOneFieldConfig{
				returnType: object,
				collection: typeConf.Collection,
			},
		)
		queryFields[object.Name()+"Nodes"] = mongoke.findManyField(
			findManyFieldConfig{
				returnType: object,
				collection: typeConf.Collection,
			},
		)

		// TODO add mutaiton fields
		mutationFields["putSome"+object.Name()] = &graphql.Field{
			Type: object,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		}
	}

	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Types:      baseSchemaConfig.Types,
			Extensions: baseSchemaConfig.Extensions,
			Query:      graphql.NewObject(graphql.ObjectConfig{Name: "Query", Fields: queryFields}),
			Mutation:   graphql.NewObject(graphql.ObjectConfig{Name: "Mutation", Fields: mutationFields}),
		},
	)
	if err != nil {
		return graphql.Schema{}, err
	}
	return schema, nil
}

/*
functions to
- WhereObject input based on fields
- ConnectionType based on nodes type
- Edge based on nodes type
- Fields enum based on fields
*/

/*
functions to
- create a many resolver, based on collection, guard
- create a one resolver, based on collection
*/

const (
	general_graphql = `
type Query {
    mongoke_version: String
}
enum Direction {
    ASC
    DESC
}
${{
''.join([f"""
input Where{scalar} {'{'}
    in: [{scalar}]
    nin: [{scalar}]
    eq: {scalar}
    neq: {scalar}
{'}'}
""" for scalar in map(str, sorted(searchables))])
}}
type PageInfo {
    startCursor: AnyScalar
    endCursor: AnyScalar
    hasNextPage: Boolean
    hasPreviousPage: Boolean
}
scalar Json
scalar ObjectId
scalar AnyScalar
`

	graphql_query = `
extend type Query {
    ${{query_name}}(
        where: ${{type_name}}Where,
    ): ${{type_name}}
    ${{query_name}}Nodes(
        where: ${{type_name}}Where, 
        cursorField: ${{type_name}}Fields, 
        first: Int, 
        last: Int, 
        after: AnyScalar, 
        before: AnyScalar,
        direction: Direction
    ): ${{type_name}}Connection!
}
type ${{type_name}}Connection {
    nodes: [${{type_name}}]!
    edges: [${{type_name}}Edge]!
    pageInfo: PageInfo!
}
type ${{type_name}}Edge {
    node: ${{type_name}}
    cursor: AnyScalar
}
input ${{type_name}}Where { 
    and: [${{type_name}}Where]
    or: [${{type_name}}Where]
    # $not: [${{type_name}}Where]
    ${{'\\n    '.join([f'{field}: Where{scalar_name}' for field, scalar_name in fields])}}
}
enum ${{type_name}}Fields {
    ${{'\\n    '.join([f'{field}' for field, _ in fields])}}
}
`

	to_one_relation = `
extend type ${{fromType}} {
    ${{relationName}}: ${{toType}}
}
`

	to_many_relation = `
extend type ${{fromType}} {
   ${{relationName}}(
       where: ${{toType}}Where, 
       cursorField: ${{toType}}Fields, 
       first: Int, 
       last: Int, 
       after: AnyScalar, 
       before: AnyScalar,
       direction: Direction
    ): ${{toType}}Connection!
}
`

	to_many_relation_boilerplate = `
type ${{toType}}Connection {
    nodes: [${{toType}}]!
    edges: [${{toType}}Edge]!
    pageInfo: PageInfo!
}
type ${{toType}}Edge {
    node: ${{toType}}
    cursor: AnyScalar
}
input ${{toType}}Where { 
    ${{'\\n    '.join([f'{field}: Where{scalar_name}' for field, scalar_name in fields])}}
}
enum ${{toType}}Fields {
    ${{'\\n    '.join([f'{field}' for field, _ in fields])}}
}
`
)

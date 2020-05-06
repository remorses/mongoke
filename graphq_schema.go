package mongoke

import (
	"github.com/graphql-go/graphql"
	tools "github.com/remorses/graphql-go-tools"
)

func generateSchema(config Config) (graphql.Schema, error) {
	typeDefs := config.schemaString
	queryFields := graphql.Fields{}
	mutationFields := graphql.Fields{}
	baseSchemaConfig, err := tools.MakeSchemaConfig(tools.ExecutableSchema{TypeDefs: []string{typeDefs}})
	if err != nil {
		return graphql.Schema{}, err
	}
	for _, gqlType := range baseSchemaConfig.Types {
		// TODO handle unions
		object, ok := gqlType.(*graphql.Object)
		if !ok {
			continue
		}
		queryFields["findOne"+object.Name()] = findOneResolver(findOneResolverConfig{resolvedType: object})
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

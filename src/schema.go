package mongoke

import (
	"errors"

	"github.com/graphql-go/graphql"
)

func (mongoke *Mongoke) generateSchema() (graphql.Schema, error) {
	queryFields := graphql.Fields{}
	mutationFields := graphql.Fields{}
	baseSchemaConfig := mongoke.schemaConfig

	// add fields
	for _, gqlType := range baseSchemaConfig.Types {
		var object graphql.Type
		switch v := gqlType.(type) {
		case *graphql.Object, *graphql.Union:
			object = v
		default:
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
		p := createFieldParams{
			returnType:  object,
			permissions: typeConf.Permissions,
			collection:  typeConf.Collection,
		}
		findOne, err := mongoke.findOneField(p)
		if err != nil {
			return graphql.Schema{}, err
		}
		queryFields[object.Name()] = findOne
		findMany, err := mongoke.findManyField(p)
		if err != nil {
			return graphql.Schema{}, err
		}
		queryFields[object.Name()+"Nodes"] = findMany

		// TODO add mutaiton fields
		mutationFields["putSome"+object.Name()] = &graphql.Field{
			Type: object,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		}
	}

	// add relations
	for _, relation := range mongoke.Config.Relations {
		if relation.Field == "" {
			return graphql.Schema{}, errors.New("relation field is empty " + relation.From)
		}
		fromType := findType(baseSchemaConfig.Types, relation.From)
		if fromType == nil {
			return graphql.Schema{}, errors.New("cannot find relation `from` type " + relation.From)
		}
		returnType := findType(baseSchemaConfig.Types, relation.To)
		if returnType == nil {
			return graphql.Schema{}, errors.New("cannot find relation `to` type " + relation.To)
		}
		returnTypeConf := mongoke.Config.getTypeConfig(relation.To)
		if returnTypeConf == nil {
			return graphql.Schema{}, errors.New("cannot find type config for relation " + relation.Field)
		}
		object, ok := fromType.(*graphql.Object)
		if !ok {
			return graphql.Schema{}, errors.New("relation return type " + fromType.Name() + " is not an object")
		}
		p := createFieldParams{
			returnType:   returnType,
			permissions:  returnTypeConf.Permissions,
			collection:   returnTypeConf.Collection,
			initialWhere: relation.Where,
			omitWhere:    true,
		}
		if relation.RelationType == "to_many" {
			field, err := mongoke.findManyField(p)
			if err != nil {
				return graphql.Schema{}, err
			}
			object.AddFieldConfig(relation.Field, field)
		} else if relation.RelationType == "to_one" {
			field, err := mongoke.findOneField(p)
			if err != nil {
				return graphql.Schema{}, err
			}
			object.AddFieldConfig(relation.Field, field)
		} else {
			return graphql.Schema{}, errors.New("relation_type must be to_many or to_one")
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

func findType(a []graphql.Type, name string) graphql.Type {
	for _, t := range a {
		if t.Name() == name {
			return t
		}
	}
	return nil
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

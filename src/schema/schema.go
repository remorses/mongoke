package schema

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/PaesslerAG/gval"
	"github.com/graphql-go/graphql"
	tools "github.com/remorses/graphql-go-tools"
	mongoke "github.com/remorses/mongoke/src"
	"github.com/remorses/mongoke/src/mongodb"
	"github.com/remorses/mongoke/src/types"
)

// MakeMongokeSchema generates the schema
func MakeMongokeSchema(config mongoke.Config) (graphql.Schema, error) {
	if config.DatabaseFunctions == nil {
		config.DatabaseFunctions = mongodb.MongodbDatabaseFunctions{}
	}
	if config.Cache == nil {
		config.Cache = make(mongoke.Map)
	}

	// TODO validate config here

	if config.Schema == "" && config.SchemaPath != "" {
		data, e := ioutil.ReadFile(config.SchemaPath)
		if e != nil {
			return graphql.Schema{}, e
		}
		config.Schema = string(data)
	}
	schemaConfig, err := makeSchemaConfig(config)
	if err != nil {
		return graphql.Schema{}, err
	}
	schema, err := GenerateSchema(config, schemaConfig)
	if err != nil {
		return schema, err
	}
	return schema, nil
}

func makeSchemaConfig(config mongoke.Config) (graphql.SchemaConfig, error) {
	resolvers := map[string]tools.Resolver{
		types.ObjectID.Name(): &tools.ScalarResolver{
			Serialize:    types.ObjectID.Serialize,
			ParseLiteral: types.ObjectID.ParseLiteral,
			ParseValue:   types.ObjectID.ParseValue,
		},
	}
	for name, typeConf := range config.Types {
		if typeConf.IsTypeOf == "" {
			continue
		}
		eval, err := gval.Full().NewEvaluable(typeConf.IsTypeOf)
		if err != nil {
			return graphql.SchemaConfig{}, errors.New("got an error parsing isTypeOf expression " + typeConf.IsTypeOf)
		}
		resolvers[name] = &tools.ObjectResolver{
			IsTypeOf: func(p graphql.IsTypeOfParams) bool {
				res, err := eval(context.Background(), mongoke.Map{
					"x":        p.Value,
					"document": p.Value,
				})
				if err != nil {
					fmt.Println("got an error evaluating expression " + typeConf.IsTypeOf)
					return false
				}
				if res == true {
					return true
				}
				return false
			},
		}
	}

	baseSchemaConfig, err := tools.MakeSchemaConfig(
		tools.ExecutableSchema{
			TypeDefs:  []string{config.Schema},
			Resolvers: resolvers,
		},
	)
	return baseSchemaConfig, err
}

func GenerateSchema(Config mongoke.Config, baseSchemaConfig graphql.SchemaConfig) (graphql.Schema, error) {
	queryFields := graphql.Fields{}
	mutationFields := graphql.Fields{}

	// add fields
	for _, gqlType := range baseSchemaConfig.Types {
		var object graphql.Type
		switch v := gqlType.(type) {
		case *graphql.Object, *graphql.Union:
			object = v
		default:
			continue
		}

		typeConf := Config.GetTypeConfig(gqlType.Name())

		if typeConf == nil || (typeConf.Exposed != nil && !*typeConf.Exposed) {
			println("ignoring not exposed type " + gqlType.Name())
			continue
		}

		if typeConf.Collection == "" {
			return graphql.Schema{}, errors.New("no collection given for type " + gqlType.Name())
		}
		p := createFieldParams{
			Config:       Config,
			returnType:   object,
			permissions:  typeConf.Permissions,
			collection:   typeConf.Collection,
			schemaConfig: baseSchemaConfig,
		}
		findOne, err := findOneField(p)
		if err != nil {
			return graphql.Schema{}, err
		}
		queryFields[object.Name()] = findOne
		findMany, err := findManyField(p)
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
	for _, relation := range Config.Relations {
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
		returnTypeConf := Config.GetTypeConfig(relation.To)
		if returnTypeConf == nil {
			return graphql.Schema{}, errors.New("cannot find type config for relation " + relation.Field)
		}
		object, ok := fromType.(*graphql.Object)
		if !ok {
			return graphql.Schema{}, errors.New("relation return type " + fromType.Name() + " is not an object")
		}
		p := createFieldParams{
			Config:       Config,
			returnType:   returnType,
			permissions:  returnTypeConf.Permissions,
			collection:   returnTypeConf.Collection,
			initialWhere: relation.Where,
			schemaConfig: baseSchemaConfig,
			omitWhere:    true,
		}
		if relation.RelationType == "to_many" {
			field, err := findManyField(p)
			if err != nil {
				return graphql.Schema{}, err
			}
			object.AddFieldConfig(relation.Field, field)
		} else if relation.RelationType == "to_one" {
			field, err := findOneField(p)
			if err != nil {
				return graphql.Schema{}, err
			}
			object.AddFieldConfig(relation.Field, field)
		} else {
			return graphql.Schema{}, errors.New("relation_type must be `to_many` or `to_one`, got " + relation.RelationType)
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

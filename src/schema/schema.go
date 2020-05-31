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
	"github.com/remorses/mongoke/src/fakedata"
	"github.com/remorses/mongoke/src/fields"
	"github.com/remorses/mongoke/src/firestore"
	"github.com/remorses/mongoke/src/mongodb"
	"github.com/remorses/mongoke/src/types"
)

// MakeMongokeSchema generates the schema
func MakeMongokeSchema(config mongoke.Config) (graphql.Schema, error) {

	if config.Cache == nil {
		config.Cache = make(mongoke.Map)
	}

	if config.Schema == "" && config.SchemaPath != "" {
		data, e := ioutil.ReadFile(config.SchemaPath)
		if e != nil {
			return graphql.Schema{}, e
		}
		config.Schema = string(data)
	}

	if config.Schema == "" && config.SchemaUrl != "" {
		data, e := mongoke.DownloadFile(config.SchemaUrl)
		if e != nil {
			return graphql.Schema{}, e
		}
		config.Schema = string(data)
	}

	if config.Schema == "" {
		return graphql.Schema{}, errors.New("missing required schema")
	}

	if config.Mongodb.Uri != "" {
		config.DatabaseUri = config.Mongodb.Uri
		config.DatabaseFunctions = &mongodb.MongodbDatabaseFunctions{}
	} else if config.Firestore.ProjectID != "" {
		config.DatabaseUri = config.Firestore.ProjectID // TODO firestore project id is not really a database uri
		config.DatabaseFunctions = &firestore.FirestoreDatabaseFunctions{}
	}

	if config.DatabaseFunctions == nil { // by default use local fake data
		println("using local fake data source")
		config.DatabaseUri = ""
		config.DatabaseFunctions = &fakedata.FakeDatabaseFunctions{Config: config}
	}

	schemaConfig, err := makeSchemaConfig(config)
	if err != nil {
		return graphql.Schema{}, err
	}
	schema, err := generateSchema(config, schemaConfig)
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
			return graphql.SchemaConfig{}, errors.New(
				"got an error parsing isTypeOf expression " + typeConf.IsTypeOf)
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

func generateSchema(Config mongoke.Config, baseSchemaConfig graphql.SchemaConfig) (graphql.Schema, error) {
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
		p := fields.CreateFieldParams{
			Config:       Config,
			ReturnType:   object,
			Permissions:  typeConf.Permissions,
			Collection:   typeConf.Collection,
			SchemaConfig: baseSchemaConfig,
		}
		findOne, err := fields.QueryTypeField(p)
		if err != nil {
			return graphql.Schema{}, err
		}
		queryFields[object.Name()] = findOne
		typeNodes, err := fields.QueryTypeNodesField(p)
		if err != nil {
			return graphql.Schema{}, err
		}
		queryFields[object.Name()+"Nodes"] = typeNodes
		typeList, err := fields.QueryTypeListField(p)
		if err != nil {
			return graphql.Schema{}, err
		}
		queryFields[object.Name()+"List"] = typeList

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
		p := fields.CreateFieldParams{
			Config:       Config,
			ReturnType:   returnType,
			Permissions:  returnTypeConf.Permissions,
			Collection:   returnTypeConf.Collection,
			InitialWhere: relation.Where,
			SchemaConfig: baseSchemaConfig,
			OmitWhere:    true,
		}
		if relation.RelationType == "to_many" {
			field, err := fields.QueryTypeNodesField(p)
			if err != nil {
				return graphql.Schema{}, err
			}
			object.AddFieldConfig(relation.Field, field)
		} else if relation.RelationType == "to_one" {
			field, err := fields.QueryTypeField(p)
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

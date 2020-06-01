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
		config.DatabaseFunctions = &mongodb.MongodbDatabaseFunctions{Config: config}
	} else if config.Firestore.ProjectID != "" {
		config.DatabaseFunctions = &firestore.FirestoreDatabaseFunctions{Config: config}
	}

	if config.DatabaseFunctions == nil { // by default use local fake data
		println("using local fake data source")
		config.DatabaseFunctions = &fakedata.FakeDatabaseFunctions{Config: config}
	}

	schema, err := generateSchema(config)
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
					"x":        p.Value, // TODO add more variables to isTypeOf expressions
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

func generateSchema(config mongoke.Config) (graphql.Schema, error) {
	baseSchemaConfig, err := makeSchemaConfig(config)
	if err != nil {
		return graphql.Schema{}, err
	}
	// add fields
	query, err := makeQuery(config, baseSchemaConfig)
	if err != nil {
		return graphql.Schema{}, err
	}
	mutation, err := makeMutation(config, baseSchemaConfig)
	if err != nil {
		return graphql.Schema{}, err
	}
	err = addRelationsFields(config, baseSchemaConfig)
	if err != nil {
		return graphql.Schema{}, err
	}

	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Types:      baseSchemaConfig.Types,
			Extensions: baseSchemaConfig.Extensions,
			Query:      query,
			Mutation:   mutation,
		},
	)
	if err != nil {
		return graphql.Schema{}, err
	}
	return schema, nil
}

func makeQuery(Config mongoke.Config, baseSchemaConfig graphql.SchemaConfig) (*graphql.Object, error) {
	queryFields := graphql.Fields{}
	for _, gqlType := range baseSchemaConfig.Types {
		var object graphql.Type
		switch v := gqlType.(type) {
		case *graphql.Object, *graphql.Union, *graphql.Interface:
			object = v
		default:
			continue
		}

		typeConf := Config.GetTypeConfig(gqlType.Name())

		if unexposedType(typeConf) {
			println("ignoring not exposed type " + gqlType.Name())
			continue
		}

		if typeConf.Collection == "" {
			return nil, errors.New("no collection given for type " + gqlType.Name())
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
			return nil, err
		}
		queryFields[object.Name()] = findOne
		typeNodes, err := fields.QueryTypeNodesField(p)
		if err != nil {
			return nil, err
		}
		queryFields[object.Name()+"Nodes"] = typeNodes
		typeList, err := fields.QueryTypeListField(p)
		if err != nil {
			return nil, err
		}
		queryFields[object.Name()+"List"] = typeList

	}
	query := graphql.NewObject(graphql.ObjectConfig{Name: "Query", Fields: queryFields})
	return query, nil
}

func unexposedType(typeConf *mongoke.TypeConfig) bool {
	return typeConf == nil || (typeConf.Exposed != nil && !*typeConf.Exposed)
}

func makeMutation(Config mongoke.Config, baseSchemaConfig graphql.SchemaConfig) (*graphql.Object, error) {
	mutationFields := graphql.Fields{}
	for _, gqlType := range baseSchemaConfig.Types {
		var object graphql.Type
		switch v := gqlType.(type) {
		case *graphql.Object, *graphql.Union, *graphql.Interface:
			object = v
		default:
			continue
		}
		typeConf := Config.GetTypeConfig(gqlType.Name())

		if unexposedType(typeConf) {
			println("ignoring not exposed type " + gqlType.Name())
			continue
		}
		p := fields.CreateFieldParams{
			Config:       Config,
			ReturnType:   object,
			Permissions:  typeConf.Permissions,
			Collection:   typeConf.Collection,
			SchemaConfig: baseSchemaConfig,
		}
		// insertOne
		insertOne, err := fields.MutationInsertOne(p)
		if err != nil {
			return nil, err
		}
		mutationFields["insert"+object.Name()] = insertOne
		// insertMany
		insertNodes, err := fields.MutationInsertNodes(p)
		if err != nil {
			return nil, err
		}
		mutationFields["insert"+object.Name()+"Nodes"] = insertNodes
		// updateOne
		updateOne, err := fields.MutationUpdateOne(p)
		if err != nil {
			return nil, err
		}
		mutationFields["update"+object.Name()] = updateOne
		updateNodes, err := fields.MutationUpdateNodes(p)
		if err != nil {
			return nil, err
		}
		mutationFields["update"+object.Name()+"Nodes"] = updateNodes
	}
	mutation := graphql.NewObject(graphql.ObjectConfig{Name: "Mutation", Fields: mutationFields})
	return mutation, nil
}

func addRelationsFields(Config mongoke.Config, baseSchemaConfig graphql.SchemaConfig) error {
	// add relations
	for _, relation := range Config.Relations {
		if relation.Field == "" {
			return errors.New("relation field is empty " + relation.From)
		}
		fromType := findType(baseSchemaConfig.Types, relation.From)
		if fromType == nil {
			return errors.New("cannot find relation `from` type " + relation.From)
		}
		returnType := findType(baseSchemaConfig.Types, relation.To)
		if returnType == nil {
			return errors.New("cannot find relation `to` type " + relation.To)
		}
		returnTypeConf := Config.GetTypeConfig(relation.To)
		if returnTypeConf == nil {
			return errors.New("cannot find type config for relation " + relation.Field)
		}
		object, ok := fromType.(*graphql.Object)
		if !ok {
			return errors.New("relation return type " + fromType.Name() + " is not an object")
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
				return err
			}
			object.AddFieldConfig(relation.Field, field)
		} else if relation.RelationType == "to_one" {
			field, err := fields.QueryTypeField(p)
			if err != nil {
				return err
			}
			object.AddFieldConfig(relation.Field, field)
		} else {
			return errors.New("relation_type must be `to_many` or `to_one`, got " + relation.RelationType)
		}
	}
	return nil
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

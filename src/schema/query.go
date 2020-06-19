package schema

import (
	"errors"

	"github.com/graphql-go/graphql"
	goke "github.com/remorses/goke/src"
	"github.com/remorses/goke/src/fields"
)

func makeQuery(Config goke.Config, baseSchemaConfig graphql.SchemaConfig) (*graphql.Object, error) {
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
			// println("ignoring not exposed type " + gqlType.Name())
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

		// findOne
		findOne, err := fields.FindOne(p)
		if err != nil {
			return nil, err
		}
		queryFields["findOne"+object.Name()] = findOne

		// findMany
		findMany, err := fields.FindMany(p)
		if err != nil {
			return nil, err
		}
		queryFields["findMany"+object.Name()] = findMany

		// relay nodes
		typeNodes, err := fields.QueryTypeNodesField(p)
		if err != nil {
			return nil, err
		}
		queryFields[object.Name()+"Nodes"] = typeNodes

	}
	query := graphql.NewObject(graphql.ObjectConfig{Name: "Query", Fields: queryFields})
	return query, nil
}

func addRelationsFields(Config goke.Config, baseSchemaConfig graphql.SchemaConfig) error {
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
			InitialWhere: relation.Where, // TODO evaluate stuff inside the where object
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
			field, err := fields.FindOne(p)
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

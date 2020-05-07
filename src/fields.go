package mongoke

import (
	"errors"

	"github.com/graphql-go/graphql"
)

const TIMEOUT_FIND = 10

type findOneFieldConfig struct {
	collection string
	returnType *graphql.Object
}

func (mongoke *Mongoke) findOneField(conf findOneFieldConfig) *graphql.Field {
	// TODO create the where argument based on the object fields
	if conf.collection == "" {
		panic(errors.New("missing collection name for " + conf.returnType.Name() + " findOneField"))
	}
	resolver := func(params graphql.ResolveParams) (interface{}, error) {
		args := params.Args
		db, err := initMongo(mongoke.mongoDbUri)
		if err != nil {
			return nil, err
		}
		document, err := findOne(db.Collection(conf.collection), args["where"])
		if err != nil {
			return nil, err
		}
		// document, err := mongoke.database.findOne()
		// prettyPrint(args)
		return document, nil
	}
	whereArg, err := mongoke.getWhereArg(conf.returnType)
	if err != nil {
		panic(err)
	}
	return &graphql.Field{
		Type: conf.returnType,
		Args: graphql.FieldConfigArgument{
			"where": &graphql.ArgumentConfig{Type: whereArg},
		},
		Resolve: resolver,
	}
}

type findManyFieldConfig struct {
	collection string
	returnType *graphql.Object
}

func (mongoke *Mongoke) findManyField(conf findManyFieldConfig) *graphql.Field {
	resolver := func(params graphql.ResolveParams) (interface{}, error) {
		args := params.Args
		db, err := initMongo(mongoke.mongoDbUri)
		if err != nil {
			return nil, err
		}
		var cursorField string = "_id"
		var direction int = ASC
		if args["cursorField"] != nil {
			cursorField = args["cursorField"].(string)
		}
		if args["direction"] != nil {
			direction = args["direction"].(int)
		}
		document, err := findMany(
			db.Collection(conf.collection),
			args["where"],
			paginationFromArgs(args),
			cursorField, // TODO how does casting work?
			direction,
		)
		if err != nil {
			return nil, err
		}
		// document, err := mongoke.database.findOne()
		// prettyPrint(args)
		return document, nil
	}
	whereArg, err := mongoke.getWhereArg(conf.returnType)
	if err != nil {
		panic(err)
	}
	connectionType, err := mongoke.getConnectionType(conf.returnType)
	if err != nil {
		panic(err)
	}
	return &graphql.Field{
		Type: connectionType,
		Args: graphql.FieldConfigArgument{
			"where":     &graphql.ArgumentConfig{Type: whereArg},
			"first":     &graphql.ArgumentConfig{Type: graphql.Int},
			"last":      &graphql.ArgumentConfig{Type: graphql.Int},
			"direction": &graphql.ArgumentConfig{Type: directionEnum},
		},
		Resolve: resolver,
	}
}

func paginationFromArgs(x interface{}) Pagination {
	args, ok := x.(map[string]interface{})
	if !ok {
		return Pagination{}
	}
	var first int
	var last int
	var after string
	var before string
	if args["first"] != nil {
		first = args["first"].(int)
	}
	if args["last"] != nil {
		last = args["last"].(int)
	}
	if args["after"] != nil {
		after = args["after"].(string)
	}
	if args["before"] != nil {
		before = args["before"].(string)
	}
	return Pagination{
		first:  first,
		last:   last,
		after:  after,
		before: before,
	}
}

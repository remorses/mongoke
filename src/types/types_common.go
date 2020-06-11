package types

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	goke "github.com/remorses/goke/src"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	True  = true
	False = false
)

var DirectionEnum = graphql.NewEnum(graphql.EnumConfig{
	Name:        "Direction",
	Description: "asc or desc",
	Values: graphql.EnumValueConfigMap{
		"ASC": &graphql.EnumValueConfig{
			Value:       goke.ASC,
			Description: "ascending",
		},
		"DESC": &graphql.EnumValueConfig{
			Value:       goke.DESC,
			Description: "Descending",
		},
	},
})

var PageInfo = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "PageInfo",
		Description: "Pagination information",
		Fields: graphql.Fields{
			"startCursor": &graphql.Field{
				Type: AnyScalar,
			},
			"endCursor": &graphql.Field{
				Type: AnyScalar,
			},
			"hasNextPage": &graphql.Field{
				Type: graphql.Boolean,
			},
			"hasPreviousPage": &graphql.Field{
				Type: graphql.Boolean,
			},
		},
	},
)

var ObjectID = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "ObjectId",
	Description: "The `bson` scalar type represents a BSON ObjectId.",
	// Serialize serializes `bson.ObjectId` to string.
	Serialize: func(value interface{}) interface{} {
		switch value := value.(type) {
		case primitive.ObjectID:
			return value.Hex()
		case *primitive.ObjectID:
			v := *value
			return v.Hex()
		default:
			return nil
		}
	},
	// ParseValue parses GraphQL variables from `string` to `bson.ObjectId`.
	ParseValue: func(value interface{}) interface{} {
		switch value := value.(type) {
		case string:
			id, _ := primitive.ObjectIDFromHex(value)
			return id
		case *string:
			id, _ := primitive.ObjectIDFromHex(*value)
			return id
		default:
			return nil
		}
	},
	// ParseLiteral parses GraphQL AST to `bson.ObjectId`.
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.StringValue:
			id, _ := primitive.ObjectIDFromHex(valueAST.Value)
			return id
		}
		return nil
	},
})

// AnyScalar json type
var AnyScalar = graphql.NewScalar(
	graphql.ScalarConfig{
		Name:        "AnyScalar",
		Description: "The `AnyScalar` scalar type represents JSON values as specified by [ECMA-404](http://www.ecma-international.org/publications/files/ECMA-ST/ECMA-404.pdf)",
		Serialize: func(value interface{}) interface{} {
			//  it seems it can already handle ObjectId and generic scalars, but how?
			return value
		},
		ParseValue: func(value interface{}) interface{} {
			return value
		},
		ParseLiteral: parseAnyScalarLiteral,
	},
)

var UpdateManyPayload = graphql.NewObject(graphql.ObjectConfig{
	Name: "UpdateNodesPayload",
	Fields: graphql.Fields{
		"count": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

func parseAnyScalarLiteral(astValue ast.Value) interface{} {
	kind := astValue.GetKind()

	switch kind {
	case kinds.StringValue:
		return astValue.GetValue()
	case kinds.BooleanValue:
		return astValue.GetValue()
	case kinds.IntValue:
		return astValue.GetValue()
	case kinds.FloatValue:
		return astValue.GetValue()
	case kinds.ObjectValue:
		obj := make(map[string]interface{})
		for _, v := range astValue.GetValue().([]*ast.ObjectField) {
			obj[v.Name.Value] = parseAnyScalarLiteral(v.Value)
		}
		return obj
	case kinds.ListValue:
		list := make([]interface{}, 0)
		for _, v := range astValue.GetValue().([]ast.Value) {
			list = append(list, parseAnyScalarLiteral(v))
		}
		return list
	default:
		return nil
	}
}

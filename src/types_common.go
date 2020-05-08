package mongoke

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var directionEnum = graphql.NewEnum(graphql.EnumConfig{
	Name:        "Direction",
	Description: "asc or desc",
	Values: graphql.EnumValueConfigMap{
		"ASC": &graphql.EnumValueConfig{
			Value:       ASC,
			Description: "ascending",
		},
		"DESC": &graphql.EnumValueConfig{
			Value:       DESC,
			Description: "Descending",
		},
	},
})

var pageInfo = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "PageInfo",
		Description: "Pagination information",
		Fields: graphql.Fields{
			"startCursor": &graphql.Field{
				Type: graphql.String, // TODO should be anyscalar
			},
			"endCursor": &graphql.Field{
				Type: graphql.String,
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

var objectID = graphql.NewScalar(graphql.ScalarConfig{
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

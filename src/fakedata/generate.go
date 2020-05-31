package fakedata

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"github.com/icrowley/fake"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FakeData struct {
	nodes          map[string]ast.Node
	scalarsMapping map[string]string
	document       *ast.Document
}

type NewFakeDataParams struct {
	typeDefs       string
	scalarsMapping map[string]string
}

const (
	dateTimeKind = "dateTimeKind"
	idKind       = "idKind"
	objectIdKind = "objectIdKind"
)

func NewFakeData(p NewFakeDataParams) (*FakeData, error) {
	document, err := parser.Parse(parser.ParseParams{
		Source: &source.Source{
			Body: []byte(p.typeDefs),
			Name: "GraphQL",
		},
	})
	if err != nil {
		return nil, err
	}

	instance := &FakeData{
		document:       document,
		scalarsMapping: p.scalarsMapping,
		nodes: map[string]ast.Node{
			"ID": &ast.ScalarDefinition{
				Kind: idKind,
			},
			"ObjectId": &ast.ScalarDefinition{
				Kind: objectIdKind,
			},
			"DateTime": &ast.ScalarDefinition{
				Kind: dateTimeKind,
			},
			"String": &ast.ScalarDefinition{
				Kind: kinds.StringValue,
			},
			"Int": &ast.ScalarDefinition{
				Kind: kinds.IntValue,
			},
			"Float": &ast.ScalarDefinition{
				Kind: kinds.FloatValue,
			},
			"Boolean": &ast.ScalarDefinition{
				Kind: kinds.BooleanValue,
			},
			// "DateTime": graphql.DateTime,
		},
	}

	for scalarName, kind := range p.scalarsMapping {
		instance.nodes[scalarName] = &ast.ScalarDefinition{
			Kind: kind,
		}
	}

	for _, definition := range document.Definitions {
		// dont overwrite existing scalars
		if instance.nodes[getNodeName(definition)] != nil {
			continue
		}
		instance.nodes[getNodeName(definition)] = definition
	}

	return instance, nil
}

func (self *FakeData) Generate(name string) (interface{}, error) {
	definition := self.getNode(name)
	if definition == nil {
		return nil, errors.New("cannot find ast node for " + name)
	}
	switch nodeKind := definition.GetKind(); nodeKind {

	// custom scalars kinds
	case dateTimeKind:
		return randomDate(), nil
	case idKind:
		return fake.DigitsN(10), nil
	case objectIdKind:
		return primitive.NewObjectID(), nil

	case kinds.StringValue:
		return fake.FirstName(), nil
	case kinds.IntValue:
		return rand.Intn(100), nil
	case kinds.FloatValue:
		return rand.Float64() * 10, nil
	case kinds.BooleanValue:
		return rand.Intn(2) == 0, nil
	case kinds.ScalarDefinition:
		return nil, nil
		node := definition.(*ast.ScalarDefinition) // TODO take scalar to type mapping from options
		res, err := self.Generate(node.Name.Value)
		if err != nil {
			println("error for " + node.Name.Value + " " + err.Error())
		}
		return res, nil
	case kinds.EnumDefinition:
		node := definition.(*ast.EnumDefinition)
		randomIndex := rand.Intn(len(node.Values))
		for i, val := range node.Values {
			if i == randomIndex {
				return val.Name.Value, nil
			}
		}

	case kinds.ObjectDefinition, kinds.UnionDefinition, kinds.InterfaceDefinition:
		dest := make(map[string]interface{})
		fields := self.getNodeFields(definition)
		for _, field := range fields {
			generated, err := self.generateField(field.Type) // TODO customize the fields type of generated data via filed names
			if err != nil {
				return nil, err
			}
			dest[field.Name.Value] = generated
			// print(field.Type.String())
		}
		return dest, nil
	default:
		return nil, errors.New("kind " + nodeKind + " was not handled")
	}
	return nil, errors.New("kind " + definition.GetKind() + " was not handled")

}

func (c FakeData) generateField(astType ast.Type) (interface{}, error) {
	switch kind := astType.GetKind(); kind {
	case kinds.List:
		var list []interface{}
		for i := 0; i < 10; i++ { // TODO list has random length
			generated, err := c.generateField(astType.(*ast.List).Type)
			if err != nil {
				return nil, err
			}
			list = append(list, generated)
		}
		return list, nil

	case kinds.NonNull:
		generated, err := c.generateField(astType.(*ast.NonNull).Type)
		if err != nil {
			return nil, err
		}
		return generated, nil

	case kinds.Named:
		t := astType.(*ast.Named)
		return c.Generate(t.Name.Value)
	}

	return nil, fmt.Errorf("invalid kind for field " + astType.GetKind())
}

func getNodeName(node ast.Node) string {
	switch node.GetKind() {
	case kinds.ObjectDefinition:
		return node.(*ast.ObjectDefinition).Name.Value
	case kinds.ScalarDefinition:
		return node.(*ast.ScalarDefinition).Name.Value
	case kinds.EnumDefinition:
		return node.(*ast.EnumDefinition).Name.Value
	case kinds.InputObjectDefinition:
		return node.(*ast.InputObjectDefinition).Name.Value
	case kinds.InterfaceDefinition:
		return node.(*ast.InterfaceDefinition).Name.Value
	case kinds.UnionDefinition:
		return node.(*ast.UnionDefinition).Name.Value
	case kinds.DirectiveDefinition:
		return node.(*ast.DirectiveDefinition).Name.Value
	}

	return ""
}

func (self FakeData) getNode(name string) ast.Node {
	if val, ok := self.nodes[name]; ok {
		return val
	}
	return nil
}

func (self FakeData) getNodeFields(object ast.Node) []*ast.FieldDefinition {
	switch v := object.(type) {
	case *ast.ObjectDefinition:
		return v.Fields
	case *ast.InterfaceDefinition:
		return v.Fields
	case *ast.UnionDefinition:
		var fields []*ast.FieldDefinition
		for _, t := range v.Types {
			node := self.getNode(t.Name.Value)
			if node == nil {
				print("cannot find union type for " + t.Name.Value)
				continue
			}
			fields = append(fields, self.getNodeFields(node)...) // TODO union should choose between one of the types
		}
		return fields
	default:
		return make([]*ast.FieldDefinition, 0)
	}
}

func randomDate() time.Time {
	min := time.Date(2018, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2025, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

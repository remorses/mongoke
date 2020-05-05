package mongoke

import (
	"errors"
	"fmt"
	"strings"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/printer"
	"github.com/graphql-go/graphql/language/source"
)

var errUnresolvedDependencies = errors.New("unresolved dependencies")

// resourcesRegistry the resourcesRegistry holds all of the types
type resourcesRegistry struct {
	resources map[string]*Resource
	document  *ast.Document
	// extensions     []graphql.Extension
	unresolvedDefs []ast.Node
	maxIterations  int
	iterations     int
}

type Resource struct {
	Name string
}

func concatenateTypeDefs(typeDefs []string) (*ast.Document, error) {
	resolvedTypes := map[string]interface{}{}
	for _, defs := range typeDefs {
		doc, err := parser.Parse(parser.ParseParams{
			Source: &source.Source{
				Body: []byte(defs),
				Name: "GraphQL",
			},
		})
		if err != nil {
			return nil, err
		}

		for _, typeDef := range doc.Definitions {
			if def := printer.Print(typeDef); def != nil {
				stringDef := strings.TrimSpace(def.(string))
				resolvedTypes[stringDef] = nil
			}
		}
	}

	typeArray := []string{}
	for def := range resolvedTypes {
		typeArray = append(typeArray, def)
	}

	return parser.Parse(parser.ParseParams{
		Source: &source.Source{
			Body: []byte(strings.Join(typeArray, "\n")),
			Name: "GraphQL",
		},
	})
}

// newRegistry creates a new registry
func newRegistry(
	typeDefs []string,
) *resourcesRegistry {
	document, err := concatenateTypeDefs(typeDefs)
	if err != nil {
		// println(err)
		panic(err)
	}
	r := &resourcesRegistry{
		resources: map[string]*Resource{},
		document:  document,
		// extensions:     make([]graphql.Extension, 0),
		unresolvedDefs: document.Definitions,
		iterations:     0,
		maxIterations:  len(document.Definitions),
	}

	return r
}

// Get gets a type from the registry
func (c *resourcesRegistry) getType(name string) (*Resource, error) {
	if val, ok := c.resources[name]; ok {
		return val, nil
	}

	if !c.willResolve(name) {
		return nil, fmt.Errorf("no definition found for type %q", name)
	}

	return nil, errUnresolvedDependencies
}

// gets the extensions for the current type
func (c *resourcesRegistry) getExtensions(name, kind string) []interface{} {
	extensions := []interface{}{}

	for _, def := range c.document.Definitions {
		if def.GetKind() == kinds.TypeExtensionDefinition {
			extDef := def.(*ast.TypeExtensionDefinition).Definition
			if extDef.Name.Value == name && extDef.GetKind() == kind {
				extensions = append(extensions, extDef)
			}
		}
	}

	return extensions
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

// determines if a node will resolve eventually or with a thunk
// false if there is no possibility
func (c *resourcesRegistry) willResolve(name string) bool {
	if _, ok := c.resources[name]; ok {
		return true
	}
	for _, n := range c.unresolvedDefs {
		if getNodeName(n) == name {
			return true
		}
	}
	return false
}

// iteratively resolves dependencies until all types are resolved
func (c *resourcesRegistry) resolveDefinitions() error {
	unresolved := []ast.Node{}

	for len(c.unresolvedDefs) > 0 && c.iterations < c.maxIterations {
		c.iterations = c.iterations + 1

		for _, definition := range c.unresolvedDefs {
			switch nodeKind := definition.GetKind(); nodeKind {
			case kinds.DirectiveDefinition:
				println("cannot use directives currently ")
				continue
			case kinds.ScalarDefinition:
				println("scalar")
			case kinds.EnumDefinition:
				println("enum")
			case kinds.InputObjectDefinition:
				println("input")
			case kinds.ObjectDefinition:
				println("object")
			case kinds.InterfaceDefinition:
				println("interface")
			case kinds.UnionDefinition:
				println("union")
			case kinds.SchemaDefinition:
				println("schema")
			}
		}

		// check if everything has been resolved
		if len(unresolved) == 0 {
			return nil
		}

		// prepare the next loop
		c.unresolvedDefs = unresolved

		if c.iterations < c.maxIterations {
			unresolved = []ast.Node{}
		}
	}

	if len(unresolved) > 0 {
		names := []string{}
		for _, n := range unresolved {
			if name := getNodeName(n); name != "" {
				names = append(names, name)
			} else {
				names = append(names, n.GetKind())
			}
		}
		return fmt.Errorf("failed to resolve all type definitions: %v", names)
	}

	return nil
}

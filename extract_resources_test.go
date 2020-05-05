package mongoke

import (
	"testing"
)

const schema1 = `
type xxx {
	name: String
	surname: Int
}
`

func TestMain(t *testing.T) {
	registry := newRegistry([]string{schema1})
	registry.resolveDefinitions()
	// println(registry)
	// println(registry.document.Definitions[0].GetKind())
}

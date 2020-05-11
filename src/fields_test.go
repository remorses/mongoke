package mongoke

import (
	"testing"

	"github.com/imdario/mergo"
)

func TestMerge(t *testing.T) {
	a := map[string]Filter{
		"x": {
			Eq: "sdf",
		},
	}
	b := map[string]Filter{
		"x": {
			Gt: "sdf",
		},
		"y": {
			Neq: "sdf",
		},
	}
	mergo.Merge(&a, b)
	t.Log(pretty(a))
}

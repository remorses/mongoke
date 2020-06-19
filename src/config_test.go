package goke_test

import (
	"testing"

	goke "github.com/remorses/goke/src"
	"github.com/remorses/goke/src/testutil"
	"github.com/stretchr/testify/require"
)

func TestMakeConfigFromYaml(t *testing.T) {
	conf, err := goke.MakeConfigFromYaml(testutil.YamlConfig)
	if err != nil {
		t.Error(err)
	}
	t.Log(testutil.Pretty(conf))
}

func TestFilterInterpolate(t *testing.T) {
	filter := goke.Filter{
		Eq: "parent.name",
		Gt: "parent.age",
	}
	res, err := filter.Interpolate(goke.Map{"parent": goke.Map{"age": 12, "name": "Mike"}})
	if err != nil {
		t.Error(err)
	}
	t.Log(testutil.Pretty(res))
	require.Equal(t, "Mike", res.Eq)
	require.Equal(t, 12, res.Gt)
}

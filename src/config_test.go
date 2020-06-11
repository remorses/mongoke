package goke_test

import (
	"testing"

	goke "github.com/remorses/goke/src"
	"github.com/remorses/goke/src/testutil"
)

func TestMakeConfigFromYaml(t *testing.T) {
	conf, err := goke.MakeConfigFromYaml(testutil.YamlConfig)
	if err != nil {
		t.Error(err)
	}
	t.Log(testutil.Pretty(conf))
}

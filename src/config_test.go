package mongoke_test

import (
	"testing"

	mongoke "github.com/remorses/mongoke/src"
	"github.com/remorses/mongoke/src/testutil"
)

func TestMakeConfigFromYaml(t *testing.T) {
	conf, err := mongoke.MakeConfigFromYaml(testutil.YamlConfig)
	if err != nil {
		t.Error(err)
	}
	t.Log(testutil.Pretty(conf))
}

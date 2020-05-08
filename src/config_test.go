package mongoke

import (
	"testing"

	"github.com/remorses/mongoke/src/testutil"
)

func TestMakeConfigFromYaml(t *testing.T) {
	conf, err := MakeConfigFromYaml(testutil.YamlConfig)
	if err != nil {
		t.Error(err)
	}
	prettyPrint(conf)

}

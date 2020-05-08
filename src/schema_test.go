package mongoke

import (
	"testing"

	"github.com/remorses/mongoke/src/testutil"
)

var config = Config{
	Schema:      testutil.UserSchema,
	DatabaseUri: testutil.MONGODB_URI,
	Types: map[string]*TypeConfig{
		"User": {Collection: "users"},
	},
}

func TestSchema(t *testing.T) {
	schema, err := MakeMongokeSchema(config, nil)
	if err != nil {
		t.Error(err)
	}
	t.Run("introspect schema", func(t *testing.T) {
		if err != nil {
			t.Error(err)
		}
		testutil.QuerySchema(t, schema, testutil.IntrospectionQuery)
		// prettyPrint("introspection", data)
	})
	t.Run("query user", func(t *testing.T) {
		data := testutil.QuerySchema(t, schema, testutil.UserQuery1)
		prettyPrint("query for user", data)
	})
}

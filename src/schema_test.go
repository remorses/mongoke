package mongoke

import (
	"os"
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

// TODO add tests

func TestSchema(t *testing.T) {
	schema, err := MakeMongokeSchema(config)
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

func TestServer(t *testing.T) {
	t.Run("server", func(t *testing.T) {
		if os.Getenv("server") == "" {
			t.Skip()
		}
		println("listening on http://localhost:8080")
		main(config)
	})
}

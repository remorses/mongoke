package mongoke

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/remorses/mongoke/src/testutil"
)

func TestSchema(t *testing.T) {
	t.Run("schema", func(t *testing.T) {
		schema, err := MakeMongokeSchema(Config{schemaString: testutil.UserSchema})
		if err != nil {
			t.Error(err)
		}
		data := testutil.QuerySchema(t, schema, testutil.UserQuery1)
		_, err = json.MarshalIndent(data, "", "   ")
		if err != nil {
			t.Error(err)
		}
		// println(string(json))
	})
	t.Run("query user", func(t *testing.T) {
		schema, err := MakeMongokeSchema(Config{schemaString: testutil.UserSchema})
		if err != nil {
			t.Error(err)
		}
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
		main(Config{
			schemaString: testutil.UserSchema,
			mongoDbUri:   "mongodb://localhost/testdb",
		})
	})
}

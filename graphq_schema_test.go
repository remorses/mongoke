package mongoke

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/remorses/mongoke/testutil"
)

var schema1 = `
type xxx {
	name: String
	surname: Int
}
`

func TestSchema(t *testing.T) {
	t.Run("schema", func(t *testing.T) {
		schema, err := generateSchema(Config{schemaString: schema1})
		if err != nil {
			t.Error(err)
		}
		res := graphql.Do(graphql.Params{Schema: schema, RequestString: testutil.IntrospectionQuery})
		json, err := json.MarshalIndent(res.Data, "", "   ")
		if err != nil {
			t.Error(err)
		}
		println(string(json))

	})
}
func TestServer(t *testing.T) {
	t.Run("server", func(t *testing.T) {
		if os.Getenv("server") == "" {
			t.Skip()
		}
		println("listening")
		main(Config{schemaString: schema1})
	})
}

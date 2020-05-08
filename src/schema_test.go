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
	databaseMock := &DatabaseInterfaceMock{
		FindManyFunc: func(p FindManyParams) (Connection, error) {
			return Connection{}, nil
		},
		FindOneFunc: func(p FindOneParams) (interface{}, error) {
			return nil, nil
		},
	}
	schema, err := MakeMongokeSchema(config, databaseMock)
	if err != nil {
		t.Error(err)
	}
	t.Run("introspect schema", func(t *testing.T) {
		if err != nil {
			t.Error(err)
		}
		testutil.QuerySchema(t, schema, testutil.IntrospectionQuery)
	})
	t.Run("query user", func(t *testing.T) {
		testutil.QuerySchema(t, schema, testutil.UserQuery1)
		prettyPrint("query for user", databaseMock.FindOneCalls())
	})
}

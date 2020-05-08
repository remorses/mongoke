package mongoke

import (
	"testing"

	"github.com/remorses/mongoke/src/testutil"
	"github.com/stretchr/testify/require"
)

var config = Config{
	Schema: `
	type User {
		name: String
		age: Int
	}
	`,
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
	t.Run("findOne query without args", func(t *testing.T) {
		query := `
		{
			User {
				name
				age
			}
		}
		`
		testutil.QuerySchema(t, schema, query)
		calls := len(databaseMock.FindOneCalls())
		require.Equal(t, 1, calls)
		where := databaseMock.FindOneCalls()[calls-1].P.Where
		// require.Equal(t, nil, where)
		t.Log(where)
	})
	t.Run("findOne query with eq", func(t *testing.T) {
		databaseMock.calls.FindOne = nil
		query := `
		{
			User(where: {name: {eq: "xxx"}}) {
				name
				age
			}
		}
		`
		testutil.QuerySchema(t, schema, query)
		calls := len(databaseMock.FindOneCalls())
		require.Equal(t, 1, calls)
		where := databaseMock.FindOneCalls()[0].P.Where
		t.Log(pretty(where))
		require.Equal(t, "xxx", where["name"].Eq)
	})
	t.Run("findMany query with first, after", func(t *testing.T) {
		databaseMock.calls.FindMany = nil
		query := `
		{
			UserNodes(first: 10, after: "xxx", where: {name: {eq: "xxx"}}) {
				nodes {
					name
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}
		`
		testutil.QuerySchema(t, schema, query)
		calls := len(databaseMock.calls.FindMany)
		require.Equal(t, 1, calls)
		p := databaseMock.calls.FindMany[0].P
		t.Log("params", pretty(p))
		require.Equal(t, 10, p.Pagination.First)
		require.Equal(t, "xxx", p.Pagination.After)

	})
}

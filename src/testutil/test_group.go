package testutil

import (
	"context"
	"testing"

	"github.com/go-test/deep"
	"github.com/graphql-go/graphql"
	mongoke "github.com/remorses/mongoke/src"
)

type NewTestGroupParams struct {
	Tests         []TestCase
	Database      mongoke.DatabaseInterface
	Collection    string
	Documents     []mongoke.Map
	DefaultSchema graphql.Schema
}

type TestCase struct {
	Name          string
	Schema        graphql.Schema
	Query         string
	Expected      mongoke.Map
	ExpectedError bool
}

func NewTestGroup(t *testing.T, p NewTestGroupParams) {
	ctx := context.Background()
	_, err := p.Database.DeleteMany(ctx, mongoke.DeleteManyParams{
		Collection: p.Collection,
	})
	if err != nil {
		t.Error(err)
	}
	for _, testCase := range p.Tests {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Log()

			// t.Log(testCase.Name)
			p.Database.InsertMany(ctx, mongoke.InsertManyParams{
				Collection: p.Collection,
				Data:       p.Documents,
			})
			schema := testCase.Schema
			if testCase.ExpectedError {
				actualErr := QuerySchemaShouldFail(t, schema, testCase.Query)
				t.Log(actualErr)
				return
			}
			res := QuerySchema(t, schema, testCase.Query)
			res = ConvertToPlainMap(res)
			expected := ConvertToPlainMap(testCase.Expected)
			t.Log("expected:", expected)
			t.Log("result:", res)
			t.Log("expected:", Pretty(expected))
			t.Log("result:", Pretty(res))
			// require.Equal(t, Pretty(res), Pretty(expected))
			if diff := deep.Equal(res, expected); diff != nil {
				t.Error(diff)
			}
			_, err := p.Database.DeleteMany(ctx, mongoke.DeleteManyParams{
				Collection: p.Collection,
			})
			if err != nil {
				t.Error(err)
			}
		})
	}

}

package testutil

import (
	"context"
	"testing"

	"github.com/go-test/deep"
	"github.com/graphql-go/graphql"
	goke "github.com/remorses/goke/src"
)

type NewTestGroupParams struct {
	Tests         []TestCase
	Database      goke.DatabaseInterface
	Collection    string
	Documents     []goke.Map
	DefaultSchema graphql.Schema
}

type TestCase struct {
	Name          string
	Schema        graphql.Schema
	Query         string
	Expected      goke.Map
	ExpectedError bool
}

func NewTestGroup(t *testing.T, p NewTestGroupParams) {
	ctx := context.Background()
	if p.Database != nil {
		_, err := p.Database.DeleteMany(ctx, goke.DeleteManyParams{
			Collection: p.Collection,
		}, nil)

		if err != nil {
			t.Error(err)
		}
	}
	for _, testCase := range p.Tests {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Log()

			// t.Log(testCase.Name)
			if p.Database != nil {
				p.Database.InsertMany(ctx, goke.InsertManyParams{
					Collection: p.Collection,
					Data:       p.Documents,
				}, nil)
			}
			schema := testCase.Schema
			if testCase.ExpectedError {
				actualErr := QuerySchemaShouldFail(t, schema, testCase.Query)
				t.Log(actualErr)
				return
			}
			res := QuerySchema(t, schema, testCase.Query)
			res = ConvertToPlainMap(res)
			expected := ConvertToPlainMap(testCase.Expected)
			// t.Log("expected:", expected)
			// t.Log("result:", res)
			t.Log("expected:", Pretty(expected))
			t.Log("result:", Pretty(res))
			// require.Equal(t, Pretty(res), Pretty(expected))
			if diff := deep.Equal(res, expected); diff != nil {
				t.Error(diff)
			}
			if p.Database != nil {
				_, err := p.Database.DeleteMany(ctx, goke.DeleteManyParams{
					Collection: p.Collection,
				}, nil)
				if err != nil {
					t.Error(err)
				}
			}
		})
	}

}

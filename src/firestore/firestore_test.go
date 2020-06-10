package firestore

import (
	"context"
	"testing"

	"cloud.google.com/go/firestore"
	mongoke "github.com/remorses/mongoke/src"
	"github.com/remorses/mongoke/src/testutil"
	"github.com/stretchr/testify/require"
)

func TestFindMany(t *testing.T) {
	collection := "users"
	ctx := context.Background()
	type Case struct {
		params   mongoke.FindManyParams
		expected []mongoke.Map
		print    bool
	}
	exampleUsers := []mongoke.Map{
		{"name": "01", "age": 1},
		{"name": "02", "age": 2},
		{"name": "03", "age": 3},
	}
	cases := map[string]Case{
		"no args": {
			params: mongoke.FindManyParams{
				Collection: collection,
				OrderBy: map[string]int{
					"age": mongoke.ASC,
				},
			},
			expected: exampleUsers,
		},
		"Limit": {
			params: mongoke.FindManyParams{
				Collection: collection,
				Limit:      2,
				OrderBy: map[string]int{
					"age": mongoke.ASC,
				},
			},
			expected: exampleUsers[:2],
		},
		"Offset": {
			params: mongoke.FindManyParams{
				Collection: collection,
				Offset:     2,
				OrderBy: map[string]int{
					"age": mongoke.ASC,
				},
			},
			expected: exampleUsers[2:],
		},
		"DESC": {
			params: mongoke.FindManyParams{
				Collection: collection,
				OrderBy:    map[string]int{"age": mongoke.DESC},
			},
			expected: mongoke.ReverseMaps(exampleUsers),
		},
		"Eq": {
			params: mongoke.FindManyParams{
				Collection: collection,
				Where: mongoke.WhereTree{
					Match: map[string]mongoke.Filter{
						"name": {
							Eq: "01",
						},
					},
				},
			},
			expected: exampleUsers[:1],
		},
		"Gt": {
			params: mongoke.FindManyParams{
				Collection: collection,
				OrderBy: map[string]int{
					"name": mongoke.ASC,
				},
				Where: mongoke.WhereTree{
					Match: map[string]mongoke.Filter{
						"name": {
							Gt: "01",
						},
					},
				},
			},
			expected: exampleUsers[1:],
		},
		"Gt and Lte": {
			params: mongoke.FindManyParams{
				Collection: collection,
				Where: mongoke.WhereTree{
					Match: map[string]mongoke.Filter{
						"age": {
							Gt:  1,
							Lte: 3,
						},
					},
				},
			},
			print:    true,
			expected: exampleUsers[1:3],
		},
	}

	// clear and insert some docs
	m := FirestoreDatabaseFunctions{
		Config: mongoke.Config{
			Firestore: mongoke.FirestoreConfig{
				ProjectID: testutil.FIRESTORE_PROJECT_ID,
			},
		},
	}
	db, err := m.Init(ctx)
	if err != nil {
		t.Error(err)
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			deleteCollection(t, db, collection)
			for _, user := range exampleUsers {
				_, _, err := db.Collection(collection).Add(ctx, user)
				if err != nil {
					t.Error(err)
				}
			}
			result, err := m.FindMany(ctx, c.params)
			if err != nil {
				t.Log(c.params.Where)
				t.Error(err)
			}
			t.Log("expected", testutil.Pretty(c.expected))
			t.Log("result", testutil.Pretty(result))
			if c.print {
				println("expected", testutil.Pretty(c.expected))
				println("result", testutil.Pretty(result))
			}
			require.Equal(t, testutil.Pretty(c.expected), testutil.Pretty(result))
		})
	}
	// deleteCollection(t, db, collection)

}

func deleteCollection(t *testing.T, db *firestore.Client, collection string) {
	ctx := context.Background()
	refs, err := db.Collection(collection).Documents(ctx).GetAll()
	if err != nil {
		t.Error(err)
	}
	for _, r := range refs {
		_, err := r.Ref.Delete(ctx)
		if err != nil {
			t.Error(err)
		}
	}
}

package firestore

import (
	"context"
	"testing"

	"cloud.google.com/go/firestore"
	goke "github.com/remorses/goke/src"
	"github.com/remorses/goke/src/testutil"
	"github.com/stretchr/testify/require"
)

func TestFindMany(t *testing.T) {
	collection := "users"
	ctx := context.Background()
	type Case struct {
		params   goke.FindManyParams
		expected []goke.Map
		print    bool
	}
	exampleUsers := []goke.Map{
		{"name": "01", "age": 1},
		{"name": "02", "age": 2},
		{"name": "03", "age": 3},
	}
	cases := map[string]Case{
		"no args": {
			params: goke.FindManyParams{
				Collection: collection,
				OrderBy: map[string]int{
					"age": goke.ASC,
				},
			},
			expected: exampleUsers,
		},
		"Limit": {
			params: goke.FindManyParams{
				Collection: collection,
				Limit:      2,
				OrderBy: map[string]int{
					"age": goke.ASC,
				},
			},
			expected: exampleUsers[:2],
		},
		"Offset": {
			params: goke.FindManyParams{
				Collection: collection,
				Offset:     2,
				OrderBy: map[string]int{
					"age": goke.ASC,
				},
			},
			expected: exampleUsers[2:],
		},
		"DESC": {
			params: goke.FindManyParams{
				Collection: collection,
				OrderBy:    map[string]int{"age": goke.DESC},
			},
			expected: goke.ReverseMaps(exampleUsers),
		},
		"Eq": {
			params: goke.FindManyParams{
				Collection: collection,
				Where: goke.WhereTree{
					Match: map[string]goke.Filter{
						"name": {
							Eq: "01",
						},
					},
				},
			},
			expected: exampleUsers[:1],
		},
		"Gt": {
			params: goke.FindManyParams{
				Collection: collection,
				OrderBy: map[string]int{
					"name": goke.ASC,
				},
				Where: goke.WhereTree{
					Match: map[string]goke.Filter{
						"name": {
							Gt: "01",
						},
					},
				},
			},
			expected: exampleUsers[1:],
		},
		"Gt and Lte": {
			params: goke.FindManyParams{
				Collection: collection,
				Where: goke.WhereTree{
					Match: map[string]goke.Filter{
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
		Config: goke.Config{
			Firestore: goke.FirestoreConfig{
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
			result, err := m.FindMany(ctx,
				c.params,
				nil,
			)
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

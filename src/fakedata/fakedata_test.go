package fakedata

import (
	"context"
	"testing"

	mongoke "github.com/remorses/mongoke/src"
	"github.com/remorses/mongoke/src/testutil"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		{"name": "01", "age": 1, "_id": primitive.NewObjectID()},
		{"name": "02", "age": 2, "_id": primitive.NewObjectID()},
		{"name": "03", "age": 3, "_id": primitive.NewObjectID()},
	}
	cases := map[string]Case{
		"no args": {
			params: mongoke.FindManyParams{
				Collection: collection,
			},
			expected: exampleUsers,
		},
		"Limit": {
			params: mongoke.FindManyParams{
				Collection: collection,
				Limit:      2,
			},
			expected: exampleUsers[:2],
		},
		"Offset": {
			params: mongoke.FindManyParams{
				Collection: collection,
				Offset:     2,
			},
			expected: exampleUsers[2:],
		},
		"DESC": {
			params: mongoke.FindManyParams{
				Collection: collection,
				OrderBy:    map[string]int{"name": mongoke.DESC},
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
		"Where and And": {
			params: mongoke.FindManyParams{
				Collection: collection,
				Where: mongoke.WhereTree{
					Match: map[string]mongoke.Filter{
						"age": {
							Gt: 1,
						},
					},
					And: []mongoke.WhereTree{
						{
							Match: map[string]mongoke.Filter{
								"age": {
									Lte: 3,
								},
							},
						},
					},
				},
			},
			print:    true,
			expected: exampleUsers[1:3],
		},
		"multiple And": {
			params: mongoke.FindManyParams{
				Collection: collection,
				Where: mongoke.WhereTree{
					And: []mongoke.WhereTree{
						{
							Match: map[string]mongoke.Filter{
								"age": {
									Gt: 1,
								},
							},
						},
						{
							Match: map[string]mongoke.Filter{
								"age": {
									Lte: 3,
								},
							},
						},
					},
				},
			},
			print:    true,
			expected: exampleUsers[1:3],
		},
		"multiple Or": {
			params: mongoke.FindManyParams{
				Collection: collection,
				Where: mongoke.WhereTree{
					Or: []mongoke.WhereTree{
						{
							Match: map[string]mongoke.Filter{
								"age": {
									Eq: 2,
								},
							},
						},
						{
							Match: map[string]mongoke.Filter{
								"age": {
									Eq: 3,
								},
							},
						},
					},
				},
			},
			print:    true,
			expected: exampleUsers[1:],
		},
	}

	// clear and insert some docs
	m := FakeDatabaseFunctions{skipDataGeneration: true}
	db, err := m.Init(ctx)
	if err != nil {
		t.Error(err)
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err = db.Collection(collection).DeleteMany(ctx, mongoke.Map{})
			if err != nil {
				t.Error(err)
			}
			for _, user := range exampleUsers {
				_, err := db.Collection(collection).InsertOne(ctx, user)
				if err != nil {
					t.Error(err)
				}
			}
			result, err := m.FindMany(ctx, c.params)
			if err != nil {
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
	_, err = db.Collection(collection).DeleteMany(ctx, mongoke.Map{})

}

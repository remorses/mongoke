package goke_test

import (
	"encoding/json"
	"testing"

	goke "github.com/remorses/goke/src"
	"github.com/remorses/goke/src/testutil"
	"github.com/stretchr/testify/require"
)

func TestMakeWhereTree(t *testing.T) {
	tests := map[string]struct {
		Json     string
		Expected goke.WhereTree
	}{
		"empty string": {
			Json: `{
				
				"field": {
					"eq": ""
				}
				
			}`,
			Expected: goke.WhereTree{
				Match: map[string]goke.Filter{
					"field": {
						Eq: "",
					},
				},
			},
		},
		"1 and": {
			Json: `{
				"and": [{
					"field": {
						"eq": "xxx"
					}
				}]
			}`,
			Expected: goke.WhereTree{
				Match: make(map[string]goke.Filter),
				And: []goke.WhereTree{
					{
						Match: map[string]goke.Filter{
							"field": {
								Eq: "xxx",
							},
						},
					},
				},
			},
		},
		"2 and": {
			Json: `{
				"and": [{
					"field1": {
						"eq": "1"
					}
				},{
					"field2": {
						"eq": "2"
					}
				}]
			}`,
			Expected: goke.WhereTree{
				Match: make(map[string]goke.Filter),
				And: []goke.WhereTree{
					{
						Match: map[string]goke.Filter{
							"field1": {
								Eq: "1",
							},
						},
					},
					{
						Match: map[string]goke.Filter{
							"field2": {
								Eq: "2",
							},
						},
					},
				},
			},
		},
		"1 or, 2 and, 1 match": {
			Json: `{
				"field": {
					"eq": 9
				},
				"or": [{
					"field": {
						"eq": "xxx"
					}
				}],
				"and": [{
					"field1": {
						"eq": "1"
					}
				},{
					"field2": {
						"eq": "2"
					}
				}]
			}`,
			Expected: goke.WhereTree{
				Match: map[string]goke.Filter{
					"field": {
						Eq: float64(9),
					},
				},
				Or: []goke.WhereTree{
					{
						Match: map[string]goke.Filter{
							"field": {
								Eq: "xxx",
							},
						},
					},
				},
				And: []goke.WhereTree{
					{
						Match: map[string]goke.Filter{
							"field1": {
								Eq: "1",
							},
						},
					},
					{
						Match: map[string]goke.Filter{
							"field2": {
								Eq: "2",
							},
						},
					},
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var x map[string]interface{}
			json.Unmarshal([]byte(test.Json), &x)
			w, err := goke.MakeWhereTree(x, nil)
			if err != nil {
				t.Error(err)
			}
			t.Log(testutil.Pretty(w))
			require.Equal(t, w, test.Expected)
		})
	}

}

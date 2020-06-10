package mongoke_test

import (
	"encoding/json"
	"testing"

	mongoke "github.com/remorses/mongoke/src"
	"github.com/remorses/mongoke/src/testutil"
	"github.com/stretchr/testify/require"
)

func TestMakeWhereTree(t *testing.T) {
	tests := map[string]struct {
		Json     string
		Expected mongoke.WhereTree
	}{
		"empty string": {
			Json: `{
				
				"field": {
					"eq": ""
				}
				
			}`,
			Expected: mongoke.WhereTree{
				Match: map[string]mongoke.Filter{
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
			Expected: mongoke.WhereTree{
				Match: make(map[string]mongoke.Filter),
				And: []mongoke.WhereTree{
					{
						Match: map[string]mongoke.Filter{
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
			Expected: mongoke.WhereTree{
				Match: make(map[string]mongoke.Filter),
				And: []mongoke.WhereTree{
					{
						Match: map[string]mongoke.Filter{
							"field1": {
								Eq: "1",
							},
						},
					},
					{
						Match: map[string]mongoke.Filter{
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
			Expected: mongoke.WhereTree{
				Match: map[string]mongoke.Filter{
					"field": {
						Eq: float64(9),
					},
				},
				Or: []mongoke.WhereTree{
					{
						Match: map[string]mongoke.Filter{
							"field": {
								Eq: "xxx",
							},
						},
					},
				},
				And: []mongoke.WhereTree{
					{
						Match: map[string]mongoke.Filter{
							"field1": {
								Eq: "1",
							},
						},
					},
					{
						Match: map[string]mongoke.Filter{
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
			w, err := mongoke.MakeWhereTree(x, nil)
			if err != nil {
				t.Error(err)
			}
			t.Log(testutil.Pretty(w))
			require.Equal(t, w, test.Expected)
		})
	}

}

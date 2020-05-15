package fields

import (
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/graphql-go/graphql"
	"github.com/imdario/mergo"
	mongoke "github.com/remorses/mongoke/src"
	"github.com/remorses/mongoke/src/testutil"
	"github.com/stretchr/testify/require"
)

func TestMerge(t *testing.T) {
	a := map[string]mongoke.Filter{
		"x": {
			Eq: "sdf",
		},
	}
	b := map[string]mongoke.Filter{
		"x": {
			Gt: "sdf",
		},
		"y": {
			Neq: "sdf",
		},
	}
	mergo.Merge(&a, b)
	t.Log(testutil.Pretty(a))
}

func TestGetJwt(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		res := getJwt(graphql.ResolveParams{
			Info: graphql.ResolveInfo{
				RootValue: mongoke.Map{
					"jwt": jwt.MapClaims{
						"email": "email",
					},
				},
			},
		})
		t.Log(testutil.Pretty(res))
		require.Equal(t, "email", res["email"])
	})
}

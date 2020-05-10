package mongoke

import (
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/require"
)

func TestExtractClaims(t *testing.T) {
	token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE1ODkwMzI0MjksImV4cCI6MTYyMDU2ODQyOSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSIsIkdpdmVuTmFtZSI6IkpvaG5ueSIsIlN1cm5hbWUiOiJSb2NrZXQiLCJFbWFpbCI6Impyb2NrZXRAZXhhbXBsZS5jb20iLCJSb2xlIjpbIk1hbmFnZXIiLCJQcm9qZWN0IEFkbWluaXN0cmF0b3IiXX0.Qt_BmT2lADjJSwKhxCureJED-RDoDDrUF2bHnYGqfOo"
	secret := "qwertyuiopasdfghjklzxcvbnm123456"
	claims, err := extractClaims(JwtConfig{Key: secret}, token)
	if err != nil {
		t.Error(err)
	}
	t.Log(pretty(claims))

	require.Equal(t, "jrocket@example.com", claims["Email"])
}
func TestGetJwt(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		res := getJwt(graphql.ResolveParams{
			Info: graphql.ResolveInfo{
				RootValue: Map{
					"jwt": jwt.MapClaims{
						"email": "email",
					},
				},
			},
		})
		t.Log(pretty(res))
		require.Equal(t, "email", res["email"])
	})
	t.Run("jwt made with extractClaims", func(t *testing.T) {
		token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE1ODkwMzI0MjksImV4cCI6MTYyMDU2ODQyOSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSIsIkdpdmVuTmFtZSI6IkpvaG5ueSIsIlN1cm5hbWUiOiJSb2NrZXQiLCJFbWFpbCI6Impyb2NrZXRAZXhhbXBsZS5jb20iLCJSb2xlIjpbIk1hbmFnZXIiLCJQcm9qZWN0IEFkbWluaXN0cmF0b3IiXX0.Qt_BmT2lADjJSwKhxCureJED-RDoDDrUF2bHnYGqfOo"
		secret := "qwertyuiopasdfghjklzxcvbnm123456"
		claims, err := extractClaims(JwtConfig{Key: secret}, token)
		if err != nil {
			t.Error(err)
		}
		res := getJwt(graphql.ResolveParams{
			Info: graphql.ResolveInfo{
				RootValue: Map{
					"jwt": claims,
				},
			},
		})
		t.Log(pretty(res))
		require.Equal(t, "jrocket@example.com", res["Email"])
	})

}

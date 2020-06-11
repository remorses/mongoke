package handler

import (
	"testing"

	goke "github.com/remorses/goke/src"
	"github.com/remorses/goke/src/testutil"
	"github.com/stretchr/testify/require"
)

func TestExtractClaims(t *testing.T) {
	token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE1ODkwMzI0MjksImV4cCI6MTYyMDU2ODQyOSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSIsIkdpdmVuTmFtZSI6IkpvaG5ueSIsIlN1cm5hbWUiOiJSb2NrZXQiLCJFbWFpbCI6Impyb2NrZXRAZXhhbXBsZS5jb20iLCJSb2xlIjpbIk1hbmFnZXIiLCJQcm9qZWN0IEFkbWluaXN0cmF0b3IiXX0.Qt_BmT2lADjJSwKhxCureJED-RDoDDrUF2bHnYGqfOo"
	secret := "qwertyuiopasdfghjklzxcvbnm123456"
	claims, err := extractClaims(goke.JwtConfig{Key: secret}, token)
	if err != nil {
		t.Error(err)
	}
	t.Log(testutil.Pretty(claims))

	require.Equal(t, "jrocket@example.com", claims["Email"])
}

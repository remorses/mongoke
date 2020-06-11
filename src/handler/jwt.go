package handler

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kofalt/go-memoize"
	jwk "github.com/lestrrat-go/jwx/jwk"
	goke "github.com/remorses/goke/src"
)

func extractClaims(config goke.JwtConfig, tokenStr string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	getKey := jwtKeyGetter(config)
	if getKey == nil {
		return jwt.MapClaims{}, nil
	}
	_, err := jwt.ParseWithClaims(tokenStr, claims, getKey)
	if err != nil {
		return jwt.MapClaims{}, err
	}
	// verify issuer and audience
	if err := validateTokenClaims(config, claims); err != nil {
		return jwt.MapClaims{}, err
	}

	return claims, nil
}

/*
On startup:

GraphQL engine will fetch the JWK and will -
first, try to parse max-age or s-maxage directive in Cache-Control header.
second, check if Expires header is present (if Cache-Control is not present), and try to parse the value as a timestamp.
If it is able to parse any of the above successfully, then it will use that parsed time to refresh/refetch the JWKs again.
If it is unable to parse, then it will not refresh the JWKs (it assumes that the provider doesn't rotate their JWKs).

While running:

While GraphQL engine is running with refreshing JWKs, in one of the refresh cycles it will -
first, try to parse max-age or s-maxage directive in Cache-Control header.
second, check if Expires header is present (if Cache-Control is not present), and try to parse the value as a timestamp.
If it is able to parse any of the above successfully, then it will use that parsed time to refresh/refetch the JWKs again. If it is unable to parse, then it will sleep for 1 minute and will start another refresh cycle.

*/

func validateTokenClaims(config goke.JwtConfig, claims jwt.MapClaims) error {
	if config.Audience != "" {
		if aud, ok := claims["aud"]; ok {
			if aud != config.Audience {
				return errors.New("jwt Audience is different")
			}
		}
	}
	if config.Issuer != "" {
		if iss, ok := claims["iss"]; ok {
			if iss != config.Issuer {
				return errors.New("jwt Issuer is different")
			}
		}
	}
	return nil
}

func validateTokenHeaders(config goke.JwtConfig, token *jwt.Token) error {
	if token.Method.Alg() == "none" {
		return errors.New("jwt algorithm none is not supported")
	}
	if config.Type != "" && token.Method.Alg() != config.Type {
		return errors.New("jwt algorithm is different")
	}
	return nil
}

// JwkCache Caches jwks for 10 minutes
var JwkCache = memoize.NewMemoizer(10*time.Minute, 10*time.Minute)

func jwtKeyGetter(config goke.JwtConfig) func(token *jwt.Token) (interface{}, error) {
	if config.Key != "" {
		return func(token *jwt.Token) (interface{}, error) {
			if err := validateTokenHeaders(config, token); err != nil {
				return nil, err
			}
			return []byte(config.Key), nil
		}
	} else if config.JwkUrl != "" {
		return func(token *jwt.Token) (interface{}, error) {
			if err := validateTokenHeaders(config, token); err != nil {
				return nil, err
			}
			result, err, _ := JwkCache.Memoize(config.JwkUrl, func() (interface{}, error) {
				fmt.Println("fetching new Jwk keys from " + config.JwkUrl)
				return jwk.FetchHTTP(config.JwkUrl)
			})

			if err != nil {
				return nil, err
			}

			set := result.(*jwk.Set)

			keyID, ok := token.Header["kid"].(string)
			if !ok {
				return nil, errors.New("expecting JWT header to have string kid")
			}

			if key := set.LookupKeyID(keyID); len(key) == 1 {
				var materialized interface{}
				err := key[0].Raw(materialized)
				if err != nil {
					return "", err
				}
				return materialized, nil
			}

			return nil, fmt.Errorf("unable to find key %q", keyID)
		}
	}
	return nil
}

package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	goke "github.com/remorses/goke/src"
	goke_schema "github.com/remorses/goke/src/schema"
)

var GRAPHIQL_PATH = "/graphiql"

// MakeGokeHandler creates an http handler
func MakeGokeHandler(config goke.Config, webUiFolder string) (http.Handler, error) {
	schema, err := goke_schema.MakeGokeSchema(config)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	api := func(w http.ResponseWriter, r *http.Request) {
		// get query
		opts := handler.NewRequestOptions(r)

		// execute graphql query
		params := graphql.Params{
			Schema:         schema,
			RequestString:  opts.Query,
			VariableValues: opts.Variables,
			OperationName:  opts.OperationName,
			Context:        ctx,
		}

		params.RootObject = getRootObject(config, r)

		result := graphql.Do(params)

		// if formatErrorFn := h.formatErrorFn; formatErrorFn != nil && len(result.Errors) > 0 {
		// 	formatted := make([]gqlerrors.FormattedError, len(result.Errors))
		// 	for i, formattedError := range result.Errors {
		// 		formatted[i] = formatErrorFn(formattedError.OriginalError())
		// 	}
		// 	result.Errors = formatted
		// }

		// render web ui if content type is html
		if !config.DisableGraphiql {
			acceptHeader := r.Header.Get("Accept")
			_, raw := r.URL.Query()["raw"]
			if !raw && !strings.Contains(acceptHeader, "application/json") && strings.Contains(acceptHeader, "text/html") {
				http.Redirect(w, r, GRAPHIQL_PATH, 302)
				return
			}
		}

		w.Header().Add("Content-Type", "application/json; charset=utf-8")

		var buff []byte

		w.WriteHeader(http.StatusOK)
		buff, _ = json.MarshalIndent(result, "", "\t")

		w.Write(buff)
	}

	r := mux.NewRouter()

	graphiql, err := makeGraphiqlHandler(webUiFolder)
	if err != nil {
		return nil, err
	}

	r.PathPrefix(GRAPHIQL_PATH).Handler(http.StripPrefix(GRAPHIQL_PATH, graphiql))
	r.Handle("/", http.HandlerFunc(api))

	return r, nil
}

func makeGraphiqlHandler(webUiFolder string) (http.Handler, error) {
	// cwd, err := os.Getwd()
	// if err != nil {
	// 	return nil, err
	// }
	assets, err := filepath.Abs(webUiFolder)
	if err != nil {
		return nil, errors.New("cannot find web ui assets in " + webUiFolder + ", " + err.Error())
	}
	// assets := path.Join(root, webUiFolder)
	// fmt.Println(assets)
	h := http.FileServer(http.Dir(assets))
	return h, nil
}

func getRootObject(config goke.Config, r *http.Request) map[string]interface{} {
	rootValue := goke.Map{
		"request": r,
	}

	// jwt token
	tknStr := r.Header.Get("Authorization")
	parts := strings.Split(tknStr, "Bearer")
	tknStr = reverseStrings(parts)[0]
	tknStr = strings.TrimSpace(tknStr)
	if tknStr == "" {
		return rootValue
	}
	claims, err := extractClaims(config.JwtConfig, tknStr)
	if err != nil {
		fmt.Println("error in handler", err)
		return rootValue
	}
	rootValue["jwt"] = claims

	// admin secret
	secret := r.Header.Get(goke.AdminSecretHeader)
	isAdmin := isAdminSecretValid(config.Admins, secret)
	rootValue["isAdmin"] = isAdmin

	return rootValue
}

func corsMiddleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		next.ServeHTTP(w, r)
	})
}

func reverseStrings(input []string) []string {
	if len(input) == 0 {
		return input
	}
	return append(reverseStrings(input[1:]), input[0])
}

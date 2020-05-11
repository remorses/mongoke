package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	mongoke "github.com/remorses/mongoke/src"

	"context"
)

var WEB_UI_PATH = "./web-ui/out/"

// MakeMongokeHandler creates an http handler
func MakeMongokeHandler(config mongoke.Config) (http.Handler, error) {
	schema, err := mongoke.MakeMongokeSchema(config)
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

		params.RootObject = getRootObject(config.JwtConfig, r)

		result := graphql.Do(params)

		// if formatErrorFn := h.formatErrorFn; formatErrorFn != nil && len(result.Errors) > 0 {
		// 	formatted := make([]gqlerrors.FormattedError, len(result.Errors))
		// 	for i, formattedError := range result.Errors {
		// 		formatted[i] = formatErrorFn(formattedError.OriginalError())
		// 	}
		// 	result.Errors = formatted
		// }

		// render web ui if content type is html
		if config.EnableGraphiql == nil || config.EnableGraphiql == &mongoke.True {
			acceptHeader := r.Header.Get("Accept")
			_, raw := r.URL.Query()["raw"]
			if !raw && !strings.Contains(acceptHeader, "application/json") && strings.Contains(acceptHeader, "text/html") {
				http.Redirect(w, r, "/graphiql", 302)
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

	graphiql, _ := makeGraphiqlHandler()

	r.PathPrefix("/graphiql").Handler(http.StripPrefix("/graphiql", graphiql))
	r.Handle("/", http.HandlerFunc(api))

	return r, nil

}

func makeGraphiqlHandler() (http.Handler, error) {
	cwd, _ := os.Getwd()
	root, _ := filepath.Abs(cwd)
	assets := path.Join(root, WEB_UI_PATH)
	fmt.Println(assets)
	h := http.FileServer(http.Dir(assets))
	return h, nil

}

func getRootObject(config mongoke.JwtConfig, r *http.Request) map[string]interface{} {
	rootValue := mongoke.Map{
		"request": r,
	}

	tknStr := r.Header.Get("Authorization")
	parts := strings.Split(tknStr, "Bearer")
	tknStr = reverseStrings(parts)[0]
	tknStr = strings.TrimSpace(tknStr)

	if tknStr == "" {
		return rootValue
	}

	claims, err := extractClaims(config, tknStr)

	if err != nil {
		fmt.Println("error in handler", err)
		return rootValue
	}

	rootValue["jwt"] = claims

	return rootValue
}

func reverseStrings(input []string) []string {
	if len(input) == 0 {
		return input
	}
	// TODO remove recursion
	return append(reverseStrings(input[1:]), input[0])
}

func corsMiddleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		next.ServeHTTP(w, r)
	})
}

package mongoke

import (
	"net/http"

	"github.com/graphql-go/handler"
)

type Config struct {
	schemaString string
}

func main(config Config) {
	schema, err := generateSchema(config)
	if err != nil {
		panic(err)
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/", h)
	http.ListenAndServe("localhost:8080", nil)
}

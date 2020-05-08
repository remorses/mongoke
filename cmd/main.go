package main

import (
	"io/ioutil"
	"log"
	"os"

	"context"
	"net/http"

	mongoke "github.com/remorses/mongoke/src"
	"github.com/urfave/cli/v2"

	"github.com/graphql-go/handler"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "path",
				Usage: "config path",
			},
		},
		Action: func(c *cli.Context) error {
			path := c.String("path")
			if path == "" {
				return cli.Exit("config path is required", 1)
			}
			data, e := ioutil.ReadFile(path)
			if e != nil {
				return cli.Exit(e, 1)
			}
			config, e := mongoke.MakeConfigFromYaml(string(data))
			if e != nil {
				return cli.Exit(e, 1)
			}
			schema, err := mongoke.MakeMongokeSchema(config)
			if err != nil {
				panic(err)
			}

			// TODO the handler should be created with jwt middleware, mongoke should wrap handler package
			h := handler.New(&handler.Config{
				Schema:   &schema,
				Pretty:   true,
				GraphiQL: true,
				RootObjectFn: func(ctx context.Context, r *http.Request) map[string]interface{} {
					return make(map[string]interface{})
				},
			})
			println("listening on http://localhost:8080")
			http.Handle("/", h)
			http.ListenAndServe("localhost:8080", nil)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

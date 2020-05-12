package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"net/http"

	mongoke "github.com/remorses/mongoke/src"
	handler "github.com/remorses/mongoke/src/handler"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "path",
				Usage: "config path",
			},
			&cli.StringFlag{
				Name:  "port",
				Value: "8080",
				Usage: "port to listen to",
			},
			&cli.StringFlag{
				Name:  "www",
				Value: "",
				Usage: "web ui assets folder",
			},
			&cli.BoolFlag{
				Name:  "localhost",
				Value: false,
				Usage: "use localhost instead of 0.0.0.0",
			},
		},
		Action: func(c *cli.Context) error {
			path := c.String("path")
			// TODO download config if MONGOKE_CONFIG_URL
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
			fmt.Println("using database_uri " + config.DatabaseUri)
			h, err := handler.MakeMongokeHandler(config, c.String("www"))
			if err != nil {
				return cli.Exit(err, 1)
			}
			http.Handle("/", h)
			port := c.String("port")
			println("listening on http://localhost:" + port)
			var host string
			if c.Bool("localhost") {
				host = "localhost:"
			} else {
				host = "0.0.0.0:"
			}
			if err = http.ListenAndServe(host+port, nil); err != nil {
				fmt.Println(err)
				return cli.Exit(err, 1)
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"io/ioutil"
	"log"
	"os"

	"net/http"

	mongoke "github.com/remorses/mongoke/src"
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
			handler, err := mongoke.MakeMongokeHandler(config, nil)
			if err != nil {
				panic(err)
			}
			port := c.String("port")
			println("listening on http://localhost:" + port)
			http.Handle("/", handler)
			http.ListenAndServe("localhost:"+port, nil)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

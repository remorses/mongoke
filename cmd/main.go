package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	goke "github.com/remorses/goke/src"
	handler "github.com/remorses/goke/src/handler"
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
				Name:  "config-url",
				Usage: "config url",
			},
			&cli.StringFlag{
				Name:  "port",
				Value: "8090",
				Usage: "port to listen to",
			},
			&cli.BoolFlag{
				Name:  "localhost",
				Value: false,
				Usage: "use localhost instead of 0.0.0.0",
			},
		},
		Action: func(c *cli.Context) error {
			path := c.String("path")
			url := c.String("config-url")
			if path == "" && url == "" {
				return cli.Exit("config path or url is required", 1)
			}
			var data string
			if path != "" {
				buff, e := ioutil.ReadFile(path)
				if e != nil {
					return cli.Exit(e, 1)
				}
				data = string(buff)
			}
			if url != "" {
				var err error
				data, err = goke.DownloadFile(url)
				if err != nil {
					return cli.Exit(err, 1)
				}
			}
			config, e := goke.MakeConfigFromYaml(data)
			if e != nil {
				return cli.Exit(e, 1)
			}
			// fmt.Println("using database_uri " + config.DatabaseUri)
			www := os.Getenv("WEB_UI_ASSETS")
			if www == "" {
				return cli.Exit(errors.New("cannot find web ui assets"), 1)
			}
			h, err := handler.MakeGokeHandler(config, www)
			if err != nil {
				return cli.Exit(err, 1)
			}
			http.Handle("/", h)
			port := c.String("port")
			println("listening on http://localhost:" + port)
			var host string
			if c.Bool("localhost") {
				host = "127.0.0.1:"
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

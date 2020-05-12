package main

import (
	"bytes"
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
				Name:  "config-url",
				Usage: "config url",
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
				data, err = downloadFile(url)
				if err != nil {
					return cli.Exit(err, 1)
				}
			}
			config, e := mongoke.MakeConfigFromYaml(data)
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

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func downloadFile(url string) (string, error) {
	// TODO i should test downloadFile
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	s := buf.String() // Does a complete copy of the bytes in the buffer.

	return s, err
}

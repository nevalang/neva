package main

import (
	"log"
	_ "net/http/pprof"
	"os"

	cli "github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name: "fbp",
		Commands: []*cli.Command{
			{
				Name:   "parse",
				Action: parse,
			},
		},
	}

	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

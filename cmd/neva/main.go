package main

import (
	"fmt"
	"os"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	app := newApp(newInterpreter(), wd)
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		return
	}
}

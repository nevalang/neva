package main

import (
	"fmt"

	"github.com/emil14/neva/internal/runtime"
)

func Filter(io runtime.IO) error {
	data, err := io.In.Port("data")
	if err != nil {
		return err
	}

	marker, err := io.In.Port("marker")
	if err != nil {
		return err
	}

	acc, err := io.Out.Port("acc")
	if err != nil {
		return err
	}

	rej, err := io.Out.Port("rej")
	if err != nil {
		return err
	}

	go func() {
		for {
			fmt.Println("FILTER: start read data")
			d := <-data
			fmt.Println("FILTER: end read data")

			fmt.Println("FILTER: start read marker")
			m := <-marker
			fmt.Println("FILTER: end read marker")

			if m.Bool() {
				fmt.Println("FILTER: end send acc")
				acc <- d
				fmt.Println("FILTER: end send acc")
				continue
			}

			fmt.Println("FILTER: end send rej")
			rej <- d
			fmt.Println("FILTER: end send rej")
		}
	}()

	return nil
}

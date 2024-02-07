package main

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

type config struct {
	Home string `env:"HOME"`
}

func parseEnv() config {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	return cfg
}

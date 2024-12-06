package db

import (
	"fmt"
	"strings"

	"github.com/firodj/textenv"
)

func init() {
	oriSetenv := textenv.SetEnv
	defer func() {
		textenv.SetEnv = oriSetenv
	}()

	textenv.SetEnv = func(k string, v string) error {
		fmt.Printf("env %s=%s\n", k, v)
		return oriSetenv(k, v)
	}

	err := textenv.LoadTestEnv("../test.env")
	if err != nil {
		if !strings.Contains(err.Error(), "no such file") {
			panic(err)
		} else {
			fmt.Println("skip load test.env")
		}
	}
}

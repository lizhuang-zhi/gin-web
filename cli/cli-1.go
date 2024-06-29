package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "example",
		Usage:   "an example cli app",
		Version: "1.0.0",
		Action: func(c *cli.Context) error {
			fmt.Println("Hello, world!")
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

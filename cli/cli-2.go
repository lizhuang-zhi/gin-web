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
		Commands: []*cli.Command{ // 子命令
			{
				Name:    "greet",
				Aliases: []string{"g"},
				Usage:   "send a greeting",
				Action: func(c *cli.Context) error {
					name := "world"
					/*
						1. go run cli-2.go 是命令,不算在c.NArg()中
						2. greet 是子命令,也算在c.NArg()中
						3. greet 后面的参数才算在c.Args()中
					*/
					if c.NArg() > 0 {
						name = c.Args().Get(0)
					}
					fmt.Printf("Hello, %s!\n", name)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
	/*
		go run cli-2.go greet  // Hello, world!
		go run cli-2.go greet leo // Hello, leo!
	*/
}

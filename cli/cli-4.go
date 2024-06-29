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
			{
				Name:    "age",
				Aliases: []string{"a"},
				Usage:   "about age",
				Action: func(c *cli.Context) error {
					fmt.Println("age action!")
					return nil
				},
				Before: func(c *cli.Context) error {
					fmt.Println("age before")
					return nil
				},
				After: func(c *cli.Context) error {
					fmt.Println("age after")
					return nil
				},
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("app, action!")
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}

	/*
		lilithgames@lilithgamesdeMacBook-Pro-5 cli % go run cli-4.go
		app, action!
		lilithgames@lilithgamesdeMacBook-Pro-5 cli % go run cli-4.go greet
		Hello, world!
		lilithgames@lilithgamesdeMacBook-Pro-5 cli % go run cli-4.go greet leo
		Hello, leo!
		lilithgames@lilithgamesdeMacBook-Pro-5 cli % go run cli-4.go age
		age before
		age action!
		age after
		lilithgames@lilithgamesdeMacBook-Pro-5 cli % go run cli-4.go greet age
		Hello, age!
		lilithgames@lilithgamesdeMacBook-Pro-5 cli % go run cli-4.go greet leo age
		Hello, leo!
	*/
}

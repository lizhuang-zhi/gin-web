package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "example",
		Usage:   "an example cli app with a before hook",
		Version: "1.0.0",
		Before: func(c *cli.Context) error {
			fmt.Println("Before hook: Initializing...")
			// 这里可以放一些初始化操作，例如：
			// - 设置环境变量
			// - 读取配置文件
			// - 初始化日志系统
			// - 检查命令行参数等
			return nil
		},
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

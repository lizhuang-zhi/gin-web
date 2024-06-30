package app

import (
	"booking-app/micro-service/cluster/common/config"
	"booking-app/micro-service/cluster/common/logger"
	"booking-app/micro-service/cluster/common/mongodb"
	"booking-app/micro-service/core/meta"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

type App struct {
	cli *cli.App // 命令行应用
	// cliContext *cli.Context

	// BeforeInit func(*App) error // 初始化前的动作
	BeforeRun func(*App) error // 运行前的动作
}

func New(name, useage string) *App {
	app := &App{}

	// 元数据
	meta.Data.Set(meta.Name, name)
	meta.Data.Set(meta.Useage, useage)

	app.cli = cli.NewApp()
	app.cli.Name = name
	app.cli.Usage = useage
	app.cli.Action = app.action
	app.cli.Before = app.init
	app.cli.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Load configuration from `FILE`",
		},
		&cli.StringFlag{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "App Version",
		},
	}

	return app
}

func (app *App) Run() {
	app.cli.Run(os.Args)
}

func (app *App) init(ctx *cli.Context) error {
	// 读取配置文件路径
	configFile := ctx.String("config")
	if configFile != "" {
		fmt.Printf("Loading configuration from: %s\n", configFile)
		err := config.InitConfig(configFile)
		if err != nil {
			return err
		}
	}

	// 读取配置文件路径
	version := ctx.String("version")
	if version != "" {
		fmt.Printf("version configuration from: %s\n", version)
		meta.Data.Set(meta.Version, version) // 版本号
	}

	// 初始化日志
	logger.NewLogger()
	defer logger.GetLogger().Sync()

	// 初始化mongodb
	mongodb.NewMongoClient()

	return nil
}

func (app *App) action(ctx *cli.Context) error {
	if app.BeforeRun != nil {
		if err := app.BeforeRun(app); err != nil {
			return err
		}
	}

	return nil
}

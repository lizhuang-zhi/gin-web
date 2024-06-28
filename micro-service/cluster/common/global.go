package common

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var (
	Config      *Server            // 服务配置
	MongoClient *mongo.Client      // MongoDB
	Logger      *zap.SugaredLogger // 日志
)

type Server struct {
	// 服务配置
	System struct {
		Server  string `mapstructure:"server"`   // 服务名称
		Version string `mapstructure:"version"`  // 服务版本
		TCPPort string `mapstructure:"tcp_port"` // TCP 端口
		RPCPort string `mapstructure:"rpc_port"` // RPC 端口
	} `mapstructure:"system"`

	// MongoDB
	MongoDB struct {
		Host string `mapstructure:"host"` // 连接地址
	} `mapstructure:"mongodb"`

	// Redis
	Redis struct {
		Host string `mapstructure:"host"` // 连接地址
	} `mapstructure:"redis"`

	// 日志配置
	Log struct {
		Level string `mapstructure:"level"` // 日志级别
		Color bool   `mapstructure:"color"` // 是否开启彩色日志
		Path  string `mapstructure:"path"`  // 是否输出日志到文件，配置空则不输出
	} `mapstructure:"log"`
}

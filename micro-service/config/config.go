package config

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
}

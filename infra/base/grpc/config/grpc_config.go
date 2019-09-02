package config

type GrpcConfig struct {
	ServerConfig `yaml:"server"`
	Clients      []ClientConfig `yaml:"clients"`
}

// 本机作为服务端，服务端ip
type ServerConfig struct {
	Addr string `yaml:"addr"`
}

// 客户端连接远程服务端，服务端信息
type ClientConfig struct {
	// 服务端名称
	AppName string `yaml:"appName"`
	// 服务端地址
	Addr string `yaml:"addr"`
}

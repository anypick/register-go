package config

type RabbitMQConfig struct {
	IpAddr   string `yaml:"ipAddr"`   // 主机ip
	Port     int    `yaml:"port"`     // 端口
	Vhost    string `yaml:"vhost"`    // vhost
	UserName string `yaml:"username"` // string
	Password string `yaml:"password"` // string
}
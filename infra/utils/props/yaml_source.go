package props

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	logrus "register-go/infra/base/log/config"
	mysql "register-go/infra/base/mysql/config"
	rabbit "register-go/infra/base/rabbitmq/config"
	"register-go/infra/base/redis/config"
	"register-go/infra/utils/common"
)

// 将yaml文件映射成结构体
type YamlSource struct {
	Application       `yaml:"application"`
	config.Redis      `yaml:"redis"`
	logrus.LogConfig  `yaml:"logrus"`
	mysql.MySqlConfig `yaml:"mysql"`
	rabbit.RabbitMQConfig `yaml:"rabbit"`
}

func NewYamlSource(filePathName string) *YamlSource {
	var (
		yamlSource = new(YamlSource)
		data       []byte
		e          error
	)
	if data, e = ioutil.ReadFile(filePathName); e != nil {
		log.Fatal(e)
		return nil
	}
	if e = yaml.Unmarshal(data, yamlSource); e != nil {
		log.Fatal(e)
		return nil
	}
	return yamlSource
}

type Application struct {
	Port string `yaml:"server.port"`
	Name string `yaml:"name"`
}

func (a Application) GetDefaultPort(defaultPort string) (string, error) {
	if common.StrIsBlank(a.Port) && common.StrIsBlank(defaultPort) {
		return "", errors.New("please setting server.port")
	}
	if common.StrIsBlank(a.Port) {
		return defaultPort, nil
	} else {
		return a.Port, nil
	}
}

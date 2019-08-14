package props

import (
	"errors"
	"github.com/anypick/register-go/infra/base/redis/config"
	"github.com/anypick/register-go/infra/utils/common"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

// 将yaml文件映射成结构体
type YamlSource struct {
	Application `yaml:"application"`
	config.Redis `yaml:"redis"`
}

func NewYamlSource(filePathName string) *YamlSource {
	var (
		yamlSource = new(YamlSource)
		data []byte
		e error
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



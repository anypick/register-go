package basees

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/sirupsen/logrus"
	"register-go/infra"
)

var esClient *elasticsearch.Client

func GetClient() *elasticsearch.Client {
	return esClient
}

type EsStarter struct {
	infra.BaseStarter
}

func (e *EsStarter) Setup(ctx infra.StarterContext) {
	var err error
	config := ctx.Yaml().EsConfig
	esConfig := elasticsearch.Config{Addresses: config.Addrs}
	if esClient, err = elasticsearch.NewClient(esConfig); err != nil {
		panic(err)
	}
	logrus.Info(esClient.Info())
}

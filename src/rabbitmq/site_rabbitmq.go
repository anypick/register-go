package rabbitmq

import (
	"register-go/infra"
	"register-go/infra/rabbitmqutil"
)

var siteRabbit *SiteRabbit

func init() {
	infra.RegisterApi(&SiteRabbit{})
}

func GetSiteRabbit() *SiteRabbit {
	return siteRabbit
}

type SiteRabbit struct {
	rabbitmqutil.RabbitOperator
}

func (r *SiteRabbit) Init() {
	siteRabbit = r
	siteRabbit.Exchange = "hmall.site.direct.exchange"
	siteRabbit.Queue = "hmall.site.direct.queue"
	siteRabbit.RoutingKey = "hmall.site.direct.queue"
}

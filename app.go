package register_go

import (
	"register-go/infra"
	"register-go/infra/base"
	"register-go/infra/base/gin"
	"register-go/infra/base/log"
	"register-go/infra/base/mysql"
	"register-go/infra/base/rabbitmq"
	"register-go/infra/base/redis"
)

func init() {
	infra.Register(&base.YamlStarter{})


	infra.Register(&baselog.LogrusStarter{})
	infra.Register(&basegin.GinStarter{})
	infra.Register(&baseredis.RedisReplicationStarter{})
	infra.Register(&basesql.MySqlStarter{})
	infra.Register(&baserabbitmq.RabbitMQStarter{})


	infra.Register(&infra.BaseInitializerStarter{})
}

package register_go

import (
	"register-go/infra"
	"register-go/infra/base"
	"register-go/infra/base/gin"
	"register-go/infra/base/grpc"
	"register-go/infra/base/log"
)

func init() {
	infra.Register(&base.YamlStarter{})

	infra.Register(&baselog.LogrusStarter{})
	infra.Register(&basegrpc.GrpcServerStarter{})
	infra.Register(&basegrpc.GrpcClientStarter{})
	infra.Register(&basegin.GinStarter{})
	//infra.Register(&baseredis.RedisReplicationStarter{})
	//infra.Register(&baseredis.RedisClusterStarter{})
	//infra.Register(&basesql.MySqlStarter{})
	//infra.Register(&baserabbitmq.RabbitMQStarter{})

	infra.Register(&infra.BaseInitializerStarter{})
}

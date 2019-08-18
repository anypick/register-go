package register_go

import (
	"register-go/infra"
	"register-go/infra/base"
	"register-go/infra/base/gin"
	"register-go/infra/base/redis"
)

func init() {
	infra.Register(&base.YamlStarter{})
	infra.Register(&basegin.GinStarter{})
	infra.Register(&baseredis.RedisReplicationStarter{})
	infra.Register(&infra.BaseInitializerStarter{})
}

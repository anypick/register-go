package baseredis

import (
	"github.com/go-redis/redis"
	"register-go/infra"
)

var redisCluster *redis.ClusterClient

func GetRedisCluster() *redis.ClusterClient {
	return redisCluster
}

type RedisClusterStarter struct {
	infra.BaseStarter
}

func (r *RedisClusterStarter) Setup(ctx infra.StarterContext) {
	config := ctx.Yaml().RedisClusterConfig
	redisCluster = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        config.Addrs,
		ReadOnly:     config.ReadOnly,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
	})
}

package baseredis

import (
	"github.com/go-redis/redis"
	"register-go/infra/base"
)

var redisCluster *redis.ClusterClient

func GetRedisCluster() *redis.ClusterClient {
	return redisCluster
}

type RedisCluster struct {
	base.YamlStarter
}






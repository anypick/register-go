package baseredis

/**
这里虽然保证了节点的写入一定落在主节点，但是一旦主节点出现故障，无法进行故障转移。
无法保证Redis高可用，这里要和redis sentinel配合使用
 */
import (
	"github.com/go-redis/redis"
	"register-go/infra"
	"register-go/infra/utils/balance"
	"register-go/infra/utils/common"
)

// 定义节点角色，主节点（只写）/从节点（只读）
type NodeRole int

const (
	MasterNode = 0
	SlaveNode  = 1
	Sentinel   = 2
	Cluster    = 3
)

var clientsMap map[NodeRole][]*redis.Client

// Sentinel信息
var (
	masterName    string
	sentinelAddrs []string
)

// 获取RedisClient
// 在这里使用自定义负载均衡算法来实现负载,也可以针对架构进行优化，例如使用HAProxy进行负载
func RedisClient(nodeRole NodeRole) *redis.Client {
	clients := clientsMap[nodeRole]
	master, slave := GetRedisBalance()
	if nodeRole == MasterNode {
		return clients[master.Bl.Next(common.NilString)]
	}
	if nodeRole == SlaveNode {
		return clients[slave.Bl.Next(common.NilString)]
	}
	if nodeRole == Sentinel {
		return sentinelClient
	}
	return nil
}

var (
	masterBalance *RedisBalance
	slaveBalance  *RedisBalance
)

type RedisBalance struct {
	Bl balance.Balance
}

func GetRedisBalance() (*RedisBalance, *RedisBalance) {
	return masterBalance, slaveBalance
}

// redis主从复制Starter
type RedisReplicationStarter struct {
	infra.BaseStarter
}

func (r *RedisReplicationStarter) Init(ctx infra.StarterContext) {
	clientsMap = make(map[NodeRole][]*redis.Client, 2)
	masterName = ctx.Yaml().RedisSentinelConfig.MasterName
	sentinelAddrs = ctx.Yaml().RedisSentinelConfig.Addrs
}

func (r *RedisReplicationStarter) Setup(context infra.StarterContext) {
	config := context.Yaml().RedisConfig
	masterClient := make([]*redis.Client, 0)
	slaveClient := make([]*redis.Client, 0)
	for _, cnf := range config {
		if cnf.ReadOnly {
			slaveClient = append(slaveClient, redis.NewClient(&redis.Options{
				Addr: cnf.Addr,
			}))
		} else {
			masterClient = append(masterClient, redis.NewClient(&redis.Options{
				Addr: cnf.Addr,
			}))
		}
	}
	clientsMap[MasterNode] = masterClient
	clientsMap[SlaveNode] = slaveClient

	// 设置Redis轮询算法
	masterBalance = &RedisBalance{Bl: &balance.RoundBalance{Size: len(masterClient)}}
	slaveBalance = &RedisBalance{Bl: &balance.RoundBalance{Size: len(slaveClient)}}
	sentinelClient = redis.NewFailoverClient(&redis.FailoverOptions{MasterName: masterName, SentinelAddrs: sentinelAddrs})
}

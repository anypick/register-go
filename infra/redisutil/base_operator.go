package redisutil

import (
	"errors"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"log"
	"register-go/infra/base/redis"
	"time"
)

// 添加pipline
func ExecutePipeline(redisOps []RedisOperation, expired time.Duration) bool {
	pipeline := baseredis.RedisClient(baseredis.MasterNode).Pipeline()
	for _, redisOp := range redisOps {
		switch redisOp.Operation {
		case Add:
			addPipeline(pipeline, redisOp, expired)
			break
		case Del:
			delPipeline(pipeline, redisOp)
			break
		default:
			break
		}
	}
	cmders, e := pipeline.Exec()
	logrus.Info(cmders)
	if e != nil {
		logrus.Error(e)
		return false
	}
	return true
}

func addPipeline(pipeline redis.Pipeliner, redisOperation RedisOperation, expired time.Duration) {
	var (
		key = redisOperation.Key
		raw = redisOperation.RawKey
		err error
	)
	switch redisOperation.DataType {
	case Strings:
		pipeline.Set(key, redisOperation.Value, expired)
		//err = statusCmd.Err()
		break
	case Hashes:
		pipeline.HSet(key, raw, redisOperation.Value)
		//err = boolCmd.Err()
		break
	case ZSets:
		pipeline.ZAdd(key, redis.Z{Score: redisOperation.Score, Member: raw})
		//err = intCmd.Err()
	case Lists:
		pipeline.RPush(key, raw)
		//err = intCmd.Err()
	default:
		err = errors.New("data type incorrect")
		break
	}
	if expired != DefaultExpired {
		pipeline.Expire(key, expired)
	}
	if err != nil {
		logrus.Error(err, key, raw)
	}
}

func delPipeline(pipeline redis.Pipeliner, redisOperation RedisOperation) {
	var err error
	switch redisOperation.DataType {
	case Strings:
		intCmd := pipeline.Del(redisOperation.Key)
		err = intCmd.Err()
		break
	case Hashes:
		intCmd := pipeline.HDel(redisOperation.Key, redisOperation.RawKey)
		err = intCmd.Err()
		break
	case ZSets:
		intCmd := pipeline.ZRem(redisOperation.Key, redis.Z{Member: redisOperation.RawKey})
		err = intCmd.Err()
	case Lists:
		intCmd := pipeline.LRem(redisOperation.Key, 0, redisOperation.RawKey)
		err = intCmd.Err()
	default:
		err = errors.New("data type incorrect")
		break
	}
	if err != nil {
		log.Fatalln(err, redisOperation.Key, redisOperation.RawKey)
	}
}

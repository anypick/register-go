package redisutil

import (
	"errors"
	"github.com/anypick/register-go/infra/base/redis"
	"github.com/go-redis/redis"
	"log"
	"time"
)

func ExecutePipeline(redisOps []RedisOperation, expired time.Duration) bool {
	pipeline := baseredis.RedisClient(baseredis.MasterNode).Pipeline()
	for _, redisOp := range redisOps {
		switch redisOp.Operation {
		case Add:
			addPipeline(pipeline, redisOp, expired)
			break;
		case Del:
			delPipeline(pipeline, redisOp)
			break;
		default:
			break;
		}
	}
	cmders, e := pipeline.Exec()
	if e != nil || len(cmders) != len(redisOps) {
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
	switch redisOperation.Type {
	case Strings:
		statusCmd := pipeline.Set(key, redisOperation.Value, expired)
		err = statusCmd.Err()
		break;
	case Hashes:
		boolCmd := pipeline.HSet(key, raw, redisOperation.Value)
		err = boolCmd.Err()
		break;
	case ZSets:
		intCmd := pipeline.ZAdd(key, redis.Z{Score: redisOperation.Score, Member: raw})
		err = intCmd.Err()
	case Lists:
		intCmd := pipeline.RPush(key, raw)
		err = intCmd.Err()
	default:
		err = errors.New("data type incorrect")
		break;
	}
	boolCmd := pipeline.Expire(key, expired)
	if boolCmd.Err() != nil {
		log.Fatalln("setting expire time error, ", err)
	}
	if err != nil {
		log.Fatalln(err, key, raw)
	}
}

func delPipeline(pipeline redis.Pipeliner, redisOperation RedisOperation) {
	var err error
	switch redisOperation.Type {
	case Strings:
		intCmd := pipeline.Del(redisOperation.Key)
		err = intCmd.Err()
		break;
	case Hashes:
		intCmd := pipeline.HDel(redisOperation.Key, redisOperation.RawKey)
		err = intCmd.Err()
		break;
	case ZSets:
		intCmd := pipeline.ZRem(redisOperation.Key, redis.Z{Member: redisOperation.RawKey})
		err = intCmd.Err()
	case Lists:
		intCmd := pipeline.LRem(redisOperation.Key, 0, redisOperation.RawKey)
		err = intCmd.Err()
	default:
		err = errors.New("data type incorrect")
		break;
	}
	if err != nil {
		log.Fatalln(err, redisOperation.Key, redisOperation.RawKey)
	}
}
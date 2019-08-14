package balance

import (
	"errors"
	"hash/crc32"
	"math/rand"
	"sync/atomic"
)

// 轮询算法
type RoundBalance struct {
	count uint32
	// 集合的长度， 必传
	Size   int
}

func (r *RoundBalance) Next(key string) int {
	if r.Size == 0 {
		panic(errors.New("size length is zero"))
	}
	return int(atomic.AddUint32(&r.count, 1)) % r.Size
}

// 随机算法
type RandomBalance struct {
	// 集合的长度， 必传
	Size int
}

func (r *RandomBalance) Next(key string) int {
	if r.Size == 0 {
		panic(errors.New("size length is zero"))
	}
	return int(rand.Uint32()) % r.Size
}

// hash算法
type HashBalance struct {
	// 集合的长度， 必传
	Size int
}

func (r *HashBalance) Next(key string) int {
	if r.Size == 0 {
		panic(errors.New("size length is zero"))
	}
	return int(crc32.ChecksumIEEE([]byte(key))) % r.Size
}

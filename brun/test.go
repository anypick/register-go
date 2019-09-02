package main

import (
	"fmt"
	"math/rand"
	"test/algorithm"
	"time"
)

// 题目说明:
// 以下为并发业务框架, 多业务 routine 并发调用流控模块接口 algorithm.FlowControl()
// 接口返回 true 表示放行, 业务继续处理
// 接口返回 false 表示拒绝, 业务放弃处理
// 流控需要实现的效果: 1秒内只允许放行 3 个业务, 不足 3 个业务则全部放行, 超过 3 个业务只放行前 3 个
//
// 要求:
// 1. main.go 不能修改
// 2. 补全 func algorithm.FlowControl() bool, 实现多业务调用的流控功能
// 3. 要求只利用 channel 实现, 不能用 sync/atomic 包

func business(id int) {
	fmt.Printf("[business] start! id=%d\n", id)
	for {
		duration := time.Millisecond * (time.Duration)(500+rand.Int()%8000)
		time.Sleep(duration)
		if ok := algorithm.FlowControl(); ok {
			fmt.Printf("[business] allow! id=%d\n", id)
			// Here our business can go on
		} else {
			fmt.Printf("[business] deny! id=%d\n", id)
		}
	}
}

func main() {
	fmt.Println("ok")

	go func() {
		for range time.NewTicker(time.Second).C {
			fmt.Println("-----------------------")
		}
	}()

	// many go-goroutine(business) call flowControl()
	for i := 0; i < 20; i++ {
		go business(i)
	}

	ch := make(chan int)
	<-ch
}

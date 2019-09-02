package main

import (
	"fmt"
	"time"
)

func main() {
	count := 0
	flag := make(chan struct{})
	for {
		select {
		case <-time.After(time.Second * 1):
			count++
			if count == 5 {
				close(flag)
			}
			fmt.Println("次数", count)
		case <-flag:
			goto OUT
		}
	}
OUT:
	fmt.Println("xxxx")
}

package main

import (
	"fmt"
	"reflect"
	"time"
)

func main() {
	var a interface{} = time.Now()
	s := reflect.ValueOf(a).Type()
	fmt.Println(s)
}

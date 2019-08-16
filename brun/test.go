package main

import (
	"fmt"
	"reflect"
	"strconv"
)

func main() {
	u := uint64(10001)
	fmt.Println(strconv.FormatUint(reflect.ValueOf(u).Uint(), 10))
}

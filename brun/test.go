package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

type SiteDto struct {
	Id         int64
	SiteId     int64
	SiteCode   string
	SiteName   string
	CreateTime time.Time
	UpdateTime time.Time
}

func main() {
	site := SiteDto{SiteId: 10001, SiteName: "上海", CreateTime: time.Now(), UpdateTime: time.Now()}
	toMap := structToMap(site)

	bytes, _ := json.Marshal(toMap)
	fmt.Println(string(bytes))
}


func structToMap(source interface{}) map[string]interface{} {
	typeOf := reflect.TypeOf(source)
	valueOf := reflect.ValueOf(source)
	count := valueOf.NumField()
	data := make(map[string]interface{}, count)
	for i := 0; i < count; i ++ {
		data[typeOf.Field(i).Name] = valueOf.Field(i).Interface()
	}
	return data
}
package common

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func StrIsNotBlank(str string) bool {
	if strings.TrimSpace(str) == "" {
		return false
	}
	return true
}

func StrIsBlank(str string) bool {
	if strings.TrimSpace(str) == "" {
		return true
	}
	return false
}

func InterfaceToStr(source interface{}) string {
	switch source.(type) {
	case string:
		return source.(string)
	case uint64, uint32, uint16, uint8, uint:
		return strconv.FormatUint(reflect.ValueOf(source).Uint(), 10)
	case int64, int32, int16, int8, int:
		return strconv.FormatInt(reflect.ValueOf(source).Int(), 10)
	case float32, float64:
		return strconv.FormatFloat(reflect.ValueOf(source).Float(), 'f', -1, 64)
	case bool:
		return strconv.FormatBool(reflect.ValueOf(source).Bool())
	case time.Time:
		return source.(time.Time).Format(TimeFormat)
	default:
		panic(errors.New("interface type cannot transform string" + reflect.ValueOf(source).Type().String()))
	}
}

func StructToMap(source interface{}) map[string]interface{} {
	typeOf := reflect.TypeOf(source)
	valueOf := reflect.ValueOf(source)
	count := valueOf.NumField()
	data := make(map[string]interface{}, count)
	for i := 0; i < count; i ++ {
		data[typeOf.Field(i).Name] = valueOf.Field(i).Interface()
	}
	return data
}

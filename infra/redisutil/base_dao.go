package redisutil

import (
	"encoding/json"
	"errors"
	"github.com/anypick/register-go/infra/base/redis"
	"github.com/anypick/register-go/infra/utils/common"
	"log"
	"strconv"
	"strings"
)

const (
	// redis key 的前缀
	Prefix = "hmall:"
	// 数据库表id字段的命名（默认）
	DefaultIdName = "id"
	// 默认多语言字段
	DefaultLangCode = "zh_CN"
	// 默认超市时间，0表示永不超时
	DefaultExpired = 0
)

type IBaseDao interface {
	Add(data map[string]interface{}, langCode string) bool
}

// 定义一个redis的key,可能需要多个属性组成，在本代码中包含如下属性,组成格式：prefix:catalog:clazz:langCode:id, 或者：prefix:catalog:clazz:langCode:tmp:fieldName:fieldValue
// prefix:catalog:clazz:langCode:id   以Id为key, struct为数据，String, Hash等数据结构
// prefix:catalog:clazz:langCode:tmp:fieldValue  以FieldValue为key, id为value, Sets数据结构
// langCode:多语言描述：zh_CN， en_GB，zh_HK...
type BaseDao struct {
	// 所属分类，
	Catalog string
	// 表对应的struct
	Clazz string
	// 数据的Id, 有些表的Id字段命名比较另类
	IdName string
	// 定义需要被查询的字段
	SelectFields []string
}

// data: 传入的数据，一定要有Id这个字段
// langCode: 多语言
// return: 是否成功
func (b *BaseDao) Add(data map[string]interface{}, langCode string) bool {
	var (
		idKey       string
		fieldKey    string
		marshalData []byte
		err         error
	)
	pipeline := baseredis.RedisClient(baseredis.MasterNode).Pipeline()
	if idKey, err = b.CreateIdKey(data[b.IdName], langCode); err != nil {
		log.Fatal(err.Error())
		return false
	}
	if marshalData, err = json.Marshal(data); err != nil {
		log.Fatal(err)
		return false
	}
	pipeline.Set(idKey, string(marshalData), DefaultExpired)

	if len(b.SelectFields) > 0 {
		for _, value := range b.SelectFields {
			if fieldKey, err = b.CreateFieldKey(data[value], value, langCode); err != nil {
				log.Fatalf(err.Error())
			}
			pipeline.SAdd(fieldKey, data[b.IdName])
		}
	}
	if _, err = pipeline.Exec(); err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

// 创建key
// id 可以是string, 也有可能是uint64,基本不会出现第三种情况
func (b *BaseDao) CreateIdKey(id interface{}, langCode string) (string, error) {
	var keyValue string
	switch id.(type) {
	case string:
		keyValue = id.(string)
		if common.StrIsBlank(keyValue) {
			return common.NilString, errors.New("cannot get id value")
		}
		break;
	case uint64:
		if id.(uint64) == 0 {
			return common.NilString, errors.New("cannot get id value")
		}
		keyValue = strconv.FormatUint(id.(uint64), 10)
		break;
	}
	if common.StrIsBlank(langCode) {
		langCode = DefaultLangCode
	}
	builder := strings.Builder{}
	builder.WriteString(Prefix)
	builder.WriteString(b.Catalog)
	builder.WriteString(":")
	builder.WriteString(b.Clazz)
	builder.WriteString(":")
	builder.WriteString(langCode)
	builder.WriteString(":")
	builder.WriteString(keyValue)
	return builder.String(), nil
}

// 创建FieldKey
func (b *BaseDao) CreateFieldKey(field interface{}, fieldName, langCode string) (string, error) {
	var (
		fieldKey    []byte
		fieldKeyStr string
		err         error
	)

	// interface{} 转string

	if common.StrIsBlank(langCode) {
		langCode = DefaultLangCode
	}
	fieldKeyStr = string(fieldKey)
	builder := strings.Builder{}
	builder.WriteString(Prefix)
	builder.WriteString(b.Catalog)
	builder.WriteString(":")
	builder.WriteString(b.Clazz)
	builder.WriteString(":")
	builder.WriteString(langCode)
	builder.WriteString(":tmp:")
	builder.WriteString(fieldName)
	builder.WriteString(":")
	builder.WriteString(fieldKeyStr)
	return builder.String(), nil
}

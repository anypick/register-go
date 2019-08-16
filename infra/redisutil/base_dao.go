package redisutil

import (
	"encoding/json"
	"errors"
	"github.com/anypick/register-go/infra/base/redis"
	"github.com/anypick/register-go/infra/utils/common"
	"github.com/go-redis/redis"
	"log"
	"strings"
	"time"
)

const (
	// redis key 的前缀
	Prefix = "hmall:"
	// 数据库表id字段的命名（默认）
	DefaultIdName = "id"
	// 默认多语言字段
	DefaultLangCode = "zh_CN"
	// 默认超时时间，0表示永不超时
	DefaultExpired = 0
)

type IBaseDao interface {
	// ========================string===================
	// string类型数据插入
	// @Param: data: 传入的数据，一定要有Id这个字段； expired: 过期时间，0表示永不过期； langCode: 多语言
	// return: 是否成功
	Add(data map[string]interface{}, expired time.Duration, langCode string) bool

	// 根据Id查询string类型数据，
	// @Param: id: 数据Id, langCode: 多语言
	// return: 返回数据
	Get(id interface{}, langCode string) (data map[string]interface{})

	// 根据属性字段查询string类型数据
	// @Param fieldName: 字段名， fieldValue：字段值， langCode：多语言
	// return: 返回查询的数据
	GetByField(fieldValue interface{}, fieldName string, langCode string) (data []map[string]interface{})

	// ======================Hash=======================

	// hash类型数据插入
	// @Param: data：传入数据，一定要有Id;expired:超时时间，0表示永不超时； langCode:多语言
	// return: 是否成功
	AddHash(data map[string]interface{}, expired time.Duration, langCode string) bool

	// hash数据类型获取
	// @Param: id： 数据Id, fieldName：字段， langCode： 多语言支持
	// return: 返回字段值，由用户转回成需要的数据类型（string->int, string->slice, string-> map ...）
	GetHash(id interface{}, fieldName, langCode string) string
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
	IdDesc FieldDescriptor
	// 定义需要被查询的字段
	SelectFields []FieldDescriptor
}

// ==================================================================================================================================
// =================================              string                ===============================================================
// ==================================================================================================================================

func (b *BaseDao) Add(data map[string]interface{}, expired time.Duration, langCode string) bool {
	var (
		idKey       string
		idValue     = data[b.IdName]
		marshalData []byte
		fieldLen    = len(b.SelectFields)
		operations  = make([]RedisOperation, fieldLen+1)
		fieldOps    []RedisOperation
		err         error
	)
	if idKey, err = b.createIdKey(idValue, langCode); err != nil {
		log.Fatal(err.Error())
		return false
	}
	if marshalData, err = json.Marshal(data); err != nil {
		log.Fatal(err)
		return false
	}
	idOperation := RedisOperation{Operation: Add, Type: Strings, Key: idKey, Value: string(marshalData)}
	if fieldOps, err = b.creatSelectedField(idValue, data, Add, langCode); err != nil {
		return false
	}
	operations[0] = idOperation
	operations = append(operations, fieldOps...)
	return ExecutePipeline(operations, DefaultExpired)
}

func (b *BaseDao) Get(id interface{}, langCode string) map[string]interface{} {
	var (
		idKey string
		err   error
	)
	if idKey, err = b.createIdKey(id, langCode); err != nil {
		log.Fatal(err)
		return nil
	}
	return b.getByKey(idKey)
}

func (b *BaseDao) GetByField(fieldValue interface{}, fieldName string, langCode string) (data []map[string]interface{}) {
	if !containerField(fieldName, b.SelectFields) {
		return nil
	}
	var (
		idKey    string
		fieldKey string
		cmders   []redis.Cmder
		err      error
	)
	if fieldKey, err = b.createFieldKey(fieldValue, fieldName, langCode); err != nil {
		log.Fatal(err)
		return nil
	}
	pipeline := baseredis.RedisClient(baseredis.SlaveNode).Pipeline()

	keys := b.getKeysByField(fieldKey)
	data = make([]map[string]interface{}, len(keys))
	for _, key := range keys {
		if idKey, err = b.createIdKey(key, langCode); err != nil {
			log.Fatal(err)
			return nil
		}
		pipeline.Get(idKey)
	}
	if cmders, err = pipeline.Exec(); err != nil {
		log.Fatal(err)
		return nil
	}
	for i, cmder := range cmders {
		stringCmd := cmder.(*redis.StringCmd)
		mapData := make(map[string]interface{})
		if err = json.Unmarshal([]byte(stringCmd.Val()), mapData); err != nil {
			log.Fatal(err)
			return nil
		}
		data[i] = mapData
	}
	return data
}

// string类型， 通过field获取id集合
func (b *BaseDao) getKeysByField(fieldKey string) []string {
	client := baseredis.RedisClient(baseredis.SlaveNode)
	return client.SMembers(fieldKey).Val()
}

// string类型，通过key获取数据
func (b *BaseDao) getByKey(key string) (data map[string]interface{}) {
	client := baseredis.RedisClient(baseredis.SlaveNode)
	result := client.Get(key)
	data = make(map[string]interface{})
	if err := json.Unmarshal([]byte(result.Val()), data); err != nil {
		log.Fatal(err)
		return nil
	}
	return data
}

// ==================================================================================================================================
// =================================              Hash                ===============================================================
// ==================================================================================================================================
func (b *BaseDao) AddHash(data map[string]interface{}, expired time.Duration, langCode string) bool {
	var (
		idKey   string
		cmders  []redis.Cmder
		err     error
		idValue = data[b.IdName]
	)
	pipeline := baseredis.RedisClient(baseredis.MasterNode).Pipeline()
	if idKey, err = b.createHashKey(langCode); err != nil {
		log.Fatalln(err)
		return false
	}
	for _, value := range data {
		pipeline.HSet(idKey, common.InterfaceToStr(idValue), value)
	}
	pipeline.Expire(idKey, expired)

	if err = b.creatSelectedField(idValue, data, langCode); err != nil {
		return false
	}

	if cmders, err = pipeline.Exec(); err != nil {
		log.Fatalln(err)
		return false
	}
	log.Println(cmders)
	return true
}

func (b *BaseDao) GetHash(id interface{}, fieldName, langCode string) string {
	var (
		idKey string
		err   error
	)
	client := baseredis.RedisClient(baseredis.SlaveNode)
	if idKey, err = b.createIdKey(id, langCode); err != nil {
		log.Fatalln(err)
		return common.NilString
	}
	return client.HGet(idKey, fieldName).Val()
}

func (b *BaseDao) GetAllHash(id interface{}, langCode string) {
	var (
		idKey string
		err   error
	)
	client := baseredis.RedisClient(baseredis.SlaveNode)
	if idKey, err = b.createIdKey(id, langCode); err != nil {
		log.Fatalln(err)
		return
	}
}

func (b *BaseDao) GetHashByField() {

}

// 创建key
// id 可以是string, 也有可能是uint64,基本不会出现第三种情况
func (b *BaseDao) createIdKey(id interface{}, langCode string) (string, error) {
	if common.StrIsBlank(langCode) {
		langCode = DefaultLangCode
	}
	idKey := common.InterfaceToStr(id)
	if common.StrIsBlank(idKey) {
		return common.NilString, errors.New("id value error")
	}
	builder := strings.Builder{}
	// TODO 使用子类去实现
	builder.WriteString(Prefix)
	builder.WriteString(b.Catalog)
	builder.WriteString(":")
	builder.WriteString(b.Clazz)
	builder.WriteString(":")
	builder.WriteString(langCode)
	builder.WriteString(":")

	builder.WriteString(idKey)
	return builder.String(), nil
}

// 创建FieldKey
func (b *BaseDao) createFieldKey(field interface{}, fieldName, langCode string) (string, error) {
	if common.StrIsBlank(langCode) {
		langCode = DefaultLangCode
	}
	fieldKey := common.InterfaceToStr(field)
	if common.StrIsBlank(fieldKey) {
		return common.NilString, errors.New("field value error")
	}
	builder := strings.Builder{}
	// TODO 使用子类去实现
	builder.WriteString(Prefix)
	builder.WriteString(b.Catalog)
	builder.WriteString(":")
	builder.WriteString(b.Clazz)
	builder.WriteString(":")
	builder.WriteString(langCode)
	builder.WriteString(":tmp:")
	builder.WriteString(fieldName)
	builder.WriteString(":")
	builder.WriteString(fieldKey)
	return builder.String(), nil
}

// 创建HashKey： key(prefix:catalog:clazz:lanCode) field(id) value
func (b *BaseDao) createHashKey(langCode string) (string, error) {
	if common.StrIsBlank(langCode) {
		langCode = DefaultLangCode
	}
	builder := strings.Builder{}
	// TODO 使用子类去实现
	builder.WriteString(Prefix)
	builder.WriteString(b.Catalog)
	builder.WriteString(":")
	builder.WriteString(b.Clazz)
	builder.WriteString(":")
	builder.WriteString(langCode)

	return builder.String(), nil
}

// 任何数据类型，添加要查询的field到redis
func (b *BaseDao) creatSelectedField(idValue interface{}, data map[string]interface{}, operation OperationType, langCode string) ([]RedisOperation, error) {
	var (
		fieldKey string
		err      error
	)
	operations := make([]RedisOperation, len(b.SelectFields))
	if len(b.SelectFields) > 0 {
		for i, value := range b.SelectFields {
			if fieldKey, err = b.createFieldKey(data[value], value, langCode); err != nil {
				log.Fatalln(err)
				return nil, err
			}
			op := RedisOperation{
				Operation: operation,
				Type:      Lists,
				Key:       fieldKey,
				RawKey:    common.InterfaceToStr(idValue),
				Score:     Zero,
			}
			operations[i] = op
		}
	}
	return operations, nil
}

// 查询是否包含fieldName
func containerField(fieldName string, selectFields []FieldDescriptor) bool {
	for i := 0; i < len(selectFields); i++ {
		if selectFields[i].FieldName == fieldName {
			return true
		}
	}
	return false
}

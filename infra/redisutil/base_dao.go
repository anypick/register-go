package redisutil

import (
	"encoding/json"
	"errors"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"log"
	"register-go/infra/base/redis"
	"register-go/infra/utils/common"
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
	// string类型数据插入
	// @Param: data: 传入的数据，一定要有Id这个字段； expired: 过期时间，0表示永不过期； langCode: 多语言
	// return: 是否成功
	Add(data map[string]interface{}, expired time.Duration, langCode string) (bool, error)

	// 根据Id查询string类型数据，
	// @Param: id: 数据Id, langCode: 多语言
	// return: 返回数据
	Get(id interface{}, langCode string) (map[string]interface{}, error)

	// 根据属性字段查询数据
	// @Param fieldName: 字段名， fieldValue：字段值， langCode：多语言， page, pageSize分页数据
	// return: 返回查询的数据
	GetByField(fieldValue interface{}, fieldName string, langCode string, page, pageSize int) ([]map[string]interface{}, error)

	// ======================Hash=======================

	// hash类型数据插入
	// @Param: data：传入数据，一定要有Id;expired:超时时间，0表示永不超时； langCode:多语言
	// return: 是否成功
	AddHash(data map[string]interface{}, expired time.Duration, langCode string) (bool, error)

	// hash数据类型获取
	// @Param: id： 数据Id, langCode： 多语言支持
	// return: 返回字段值，由用户转回成需要的数据类型（string->int, string->slice, string-> map ...）
	GetHash(id interface{}, langCode string) (string, error)

	GetAllHash(page, pageSize int, langCode string) ([]map[string]interface{}, error)
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

func (b *BaseDao) Add(data map[string]interface{}, expired time.Duration, langCode string) (bool, error) {
	var (
		idKey       string
		idValue     = data[b.IdDesc.FieldName]
		marshalData []byte
		operations  = make([]RedisOperation, 0)
		fieldOps    []RedisOperation
		err         error
	)
	if idKey, err = b.createIdKey(idValue, langCode); err != nil {
		logrus.Error(err.Error())
		return false, err
	}
	if marshalData, err = json.Marshal(data); err != nil {
		logrus.Error(err)
		return false, err
	}
	idOperation := RedisOperation{Operation: Add, DataType: Strings, Key: idKey, Value: string(marshalData)}
	if fieldOps, err = b.creatSelectedField(idValue, data, Add, langCode); err != nil {
		return false, err
	}
	operations = append(operations, idOperation)
	operations = append(operations, fieldOps...)
	return ExecutePipeline(operations, DefaultExpired), err
}

func (b *BaseDao) Get(id interface{}, langCode string) (map[string]interface{}, error) {
	var (
		idKey  string
		result *redis.StringCmd
		data   map[string]interface{}
		err    error
	)
	if idKey, err = b.createIdKey(id, langCode); err != nil {
		log.Fatal(err)
		return nil, err
	}

	client := baseredis.RedisClient(baseredis.SlaveNode)
	if result = client.Get(idKey); result.Err() != nil {
		logrus.Error(result.Err())
		return nil, result.Err()
	}
	if common.StrIsBlank(result.String()) {
		logrus.Info("数据不存在")
		return nil, nil
	}
	data = make(map[string]interface{})
	if err := json.Unmarshal([]byte(result.Val()), &data); err != nil {
		logrus.Error(err)
		return nil, err
	}
	return data, err
}

func (b *BaseDao) GetByField(fieldValue interface{}, fieldName string, langCode string, page, pageSize int) ([]map[string]interface{}, error) {
	var (
		idKey  string
		cmders []redis.Cmder
		data   []map[string]interface{}
		err    error
	)
	pipeline := baseredis.RedisClient(baseredis.SlaveNode).Pipeline()
	keys := b.getKeysByField(fieldName, common.InterfaceToStr(fieldValue), langCode, page, pageSize)
	if len(keys) == 0 {
		logrus.Warn("查询的数据不存在")
		return nil, nil
	}
	data = make([]map[string]interface{}, len(keys))
	for _, key := range keys {
		if idKey, err = b.createIdKey(key, langCode); err != nil {
			logrus.Error(err)
			return nil, err
		}
		pipeline.Get(idKey)
	}
	if cmders, err = pipeline.Exec(); err != nil {
		logrus.Error(err)
		return nil, err
	}
	for i, cmder := range cmders {
		stringCmd := cmder.(*redis.StringCmd)
		mapData := make(map[string]interface{})
		if err = json.Unmarshal([]byte(stringCmd.Val()), &mapData); err != nil {
			logrus.Error(err)
			return nil, err
		}
		data[i] = mapData
	}
	return data, err
}

// 通过field获取key
func (b *BaseDao) getKeysByField(fieldName, fieldValue, langCode string, page, pageSize int) []string {
	var (
		isContainer     bool
		fieldDescriptor FieldDescriptor
		key             string
		client          = baseredis.RedisClient(baseredis.SlaveNode)
		keys            = make([]string, 0)
	)
	// 判断field类型
	if isContainer, fieldDescriptor = containerField(fieldName, b.SelectFields); isContainer {
		switch fieldDescriptor.FieldType {
		case TypeEq:
			key, _ = b.createHashFieldKey(fieldName, langCode)
			stringCmd := client.HGet(key, fieldValue)
			if stringCmd.Err() != nil {
				logrus.Warn(stringCmd.Err())
			} else {
				keys = append(keys, stringCmd.Val())
			}
			break
		case TypeMatch:
			key, _ = b.createFieldKey(fieldValue, fieldName, langCode)
			start, stop := countPage(page, pageSize)
			stringSliceCmd := client.ZRange(key, start, stop)
			if stringSliceCmd.Err() != nil {
				logrus.Warn(stringSliceCmd.Err())
			} else {
				keys = append(keys, stringSliceCmd.Val()...)
			}
			break
		case TypeRange:
			key, _ = b.createRangeKey(fieldName, langCode)
			start, stop := countPage(page, pageSize)
			stringSliceCmd := client.ZRange(key, start, stop)
			if stringSliceCmd.Err() != nil {
				logrus.Warn(stringSliceCmd.Err())
			} else {
				keys = append(keys, stringSliceCmd.Val()...)
			}
			break
		default:
			logrus.Warn("查询field type不存在")
			return keys
		}
		return keys
	}
	logrus.Warn("查询field不存在")
	return keys
}

// ==================================================================================================================================
// =================================              Hash                ===============================================================
// ==================================================================================================================================
func (b *BaseDao) AddHash(data map[string]interface{}, expired time.Duration, langCode string) (bool, error) {
	var (
		idKey        string
		redisOps     = make([]RedisOperation, 0)
		fieldOps     []RedisOperation
		marshalValue []byte
		err          error
		idValue      = data[b.IdDesc.FieldName]
	)
	if idKey, err = b.createHashKey(langCode); err != nil {
		log.Fatalln(err)
		return false, err
	}

	if marshalValue, err = json.Marshal(data); err != nil {
		log.Fatalln(err)
		return false, err
	}
	idOps := RedisOperation{Operation: Add, DataType: Hashes, Key: idKey, RawKey: common.InterfaceToStr(idValue), Value: string(marshalValue)}
	redisOps = append(redisOps, idOps)

	if fieldOps, err = b.creatSelectedField(idValue, data, Add, langCode); err != nil {
		log.Fatalln(err)
		return false, err
	}
	redisOps = append(redisOps, fieldOps...)
	return ExecutePipeline(redisOps, expired), err
}

func (b *BaseDao) GetHash(id interface{}, langCode string) (string, error) {
	var (
		hashKey string
		err     error
	)
	client := baseredis.RedisClient(baseredis.SlaveNode)
	if hashKey, err = b.createHashKey(langCode); err != nil {
		log.Fatalln(err)
		return common.NilString, err
	}
	return client.HGet(hashKey, common.InterfaceToStr(id)).Val(), err
}

func (b *BaseDao) GetAllHash(page, pageSize int, langCode string) ([]map[string]interface{}, error) {
	var (
		hashKey  string
		client   = baseredis.RedisClient(baseredis.SlaveNode)
		pipeline = client.Pipeline()
		cmders   []redis.Cmder
		datas    = make([]map[string]interface{}, 0)
		err      error
	)
	if hashKey, err = b.createHashKey(langCode); err != nil {
		return nil, err
	}
	keys := client.HKeys(hashKey).Val()
	if keys == nil || len(keys) == 0 {
		err = errors.New("hash table is nil")
		return nil, err
	}
	start, end := pageCount(len(keys), page, pageSize)
	curPageKeys := keys[start:end]
	for _, idKey := range curPageKeys {
		pipeline.HGet(hashKey, idKey)
	}
	if cmders, err = pipeline.Exec(); err != nil {
		return nil, err
	}

	for _, value := range cmders {
		var data = make(map[string]interface{})
		val := value.(*redis.StringCmd).Val()
		if err := json.Unmarshal([]byte(val), &data); err != nil {
			logrus.Error("Unmarshal error, ", err)
			return nil, err
		}
		datas = append(datas, data)
	}
	return datas, nil
}

// =========================================================================================
// ======================   创建各种数据类型的Key   ===========================================
//===========================================================================================

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
	prefix := b.createKey(langCode)
	return prefix + ":" + idKey, nil
}

// 创建FieldKey
func (b *BaseDao) createFieldKey(fieldValue interface{}, fieldName, langCode string) (string, error) {
	if common.StrIsBlank(langCode) {
		langCode = DefaultLangCode
	}
	fieldValueKey := common.InterfaceToStr(fieldValue)
	if common.StrIsBlank(fieldValueKey) {
		return common.NilString, errors.New("field value error")
	}
	prefix := b.createKey(langCode)
	return prefix + ":" + fieldName + ":" + fieldValueKey, nil
}

// feild存储类型为hash，创建hashFieldkey
func (b *BaseDao) createHashFieldKey(fieldName, langCode string) (string, error) {
	if common.StrIsBlank(langCode) {
		langCode = DefaultLangCode
	}
	prefix := b.createKey(langCode)
	return prefix + ":" + fieldName, nil
}

// 创建HashKey：
func (b *BaseDao) createHashKey(langCode string) (string, error) {
	if common.StrIsBlank(langCode) {
		langCode = DefaultLangCode
	}
	prefix := b.createKey(langCode)
	return prefix, nil
}

// 创建排序字段的key
func (b *BaseDao) createRangeKey(fieldName, langCode string) (string, error) {
	if common.StrIsBlank(langCode) {
		langCode = DefaultLangCode
	}
	prefix := b.createKey(langCode)
	return prefix + ":" + fieldName, nil
}

// 创建Key的前缀
func (b *BaseDao) createKey(langCode string) string {
	if common.StrIsBlank(langCode) {
		langCode = DefaultLangCode
	}
	builder := strings.Builder{}
	builder.WriteString(Prefix)
	builder.WriteString(b.Catalog)
	builder.WriteString(":")
	builder.WriteString(langCode)
	builder.WriteString(":")
	builder.WriteString(b.Clazz)
	return builder.String()

}

// 任何数据类型，添加要查询的field到redis
func (b *BaseDao) creatSelectedField(idValue interface{}, data map[string]interface{}, operation OperationType, langCode string) ([]RedisOperation, error) {
	var (
		fieldKey string
		dataType = ZSets
		score    = Zero
		rawKey   = common.InterfaceToStr(idValue)
		err      error
	)
	operations := make([]RedisOperation, len(b.SelectFields))
	if len(b.SelectFields) > 0 {
		for i, value := range b.SelectFields {
			fieldValue := data[value.FieldName]
			switch value.FieldType {
			// feild和Id一一对应
			case TypeEq:
				// 使用hash数据类型
				dataType = Hashes
				if fieldKey, err = b.createHashFieldKey(value.FieldName, langCode); err != nil {
					logrus.Error(err)
					return nil, err
				}
				rawKey = common.InterfaceToStr(fieldValue)
				break
			// 需要针对Field进行排序，排序规则根据Score
			case TypeRange:
				if fieldKey, err = b.createRangeKey(value.FieldName, langCode); err != nil {
					log.Fatalln(err)
					return nil, err
				}
				dataType = ZSets
				rawKey = common.InterfaceToStr(idValue)
				score = countScore(fieldValue)
				break
			default:
				if fieldKey, err = b.createFieldKey(fieldValue, value.FieldName, langCode); err != nil {
					log.Fatalln(err)
					return nil, err
				}
				dataType = ZSets
				rawKey = common.InterfaceToStr(idValue)
				break
			}
			op := RedisOperation{
				Operation: operation,
				DataType:  dataType,
				Key:       fieldKey,
				RawKey:    rawKey,
				Score:     score,
				Value:     common.InterfaceToStr(idValue),
			}
			operations[i] = op
		}
	}
	return operations, nil
}

// 创建字段描述
func (b *BaseDao) CreateFieldDescriptor(fieldName string, fieldType FieldType) {
	if b.SelectFields == nil {
		b.SelectFields = make([]FieldDescriptor, 0)
	}
	b.SelectFields = append(b.SelectFields, FieldDescriptor{FieldName: fieldName, FieldType: fieldType})
}

// 查询是否包含fieldName
func containerField(fieldName string, selectFields []FieldDescriptor) (bool, FieldDescriptor) {
	if selectFields == nil || len(selectFields) == 0 {
		return false, FieldDescriptor{}
	}
	for i := 0; i < len(selectFields); i++ {
		if selectFields[i].FieldName == fieldName {
			return true, selectFields[i]
		}
	}
	return false, FieldDescriptor{}
}

// 根据传入字段进行计算分数
func countScore(v interface{}) float64 {
	switch v.(type) {
	case time.Time:
		return float64(v.(time.Time).Unix())
	default:
		break
	}
	return 0
}

func countPage(page, pageSize int) (int64, int64) {
	return int64((page - 1) * pageSize), int64(page*pageSize - 1)
}

// 切片的分页截取方法
func pageCount(total, page, pageSize int) (start int, end int) {
	pageNum := (total + pageSize - 1) / pageSize
	if page > pageNum {
		page = pageNum
	}
	start = (page - 1) * pageSize
	if page*pageSize > total {
		end = total
	} else {
		end = page * pageSize
	}
	return
}

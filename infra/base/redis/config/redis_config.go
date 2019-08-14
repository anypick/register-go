package config

/**
定义Redis的配置信息，映射yaml文件中的配置
 */

import (
	"errors"
	"github.com/anypick/register-go/infra/utils/common"
	"reflect"
	"time"
)

type Redis struct {
	RedisConfig []RedisConfig `yaml:"redisConfig,flow"`
}

// 定义Redis需要的配置
type RedisConfig struct {
	Addr         string        `yaml:"addr,omitempty"`
	Password     string        `yaml:"password,omitempty"`
	DB           int           `yaml:"db,omitempty"`
	MaxRetries   int           `yaml:"maxRetries,omitempty"`
	PoolSize     int           `yaml:"poolSize,omitempty"`
	MinIdleConns int           `yaml:"minIdleConns,omitempty"`
	MaxConnAge   time.Duration `yaml:"maxConnAge,omitempty"`
	ReadOnly     bool          `yaml:"readOnly,omitempty"`
}

// key为字段名称，大小写要一致
func (t RedisConfig) GetString(key string) (string, error) {
	valueOf := reflect.ValueOf(t)
	value := valueOf.FieldByName(key).Interface().(string)
	if common.StrIsBlank(value) {
		return "", errors.New("please setting" + key)
	}
	return value, nil
}

func (t RedisConfig) GetInt(key string) (int, error) {
	valueOf := reflect.ValueOf(t)
	value := valueOf.FieldByName(key).Interface().(int)
	if value == 0 {
		return value, errors.New("please setting" + key)
	}
	return value, nil
}

func (t RedisConfig) GetBool(key string) (bool, error) {
	valueOf := reflect.ValueOf(t)
	value := valueOf.FieldByName(key).Interface().(bool)
	return value, nil
}

func (t RedisConfig) GetTime(key string) (time.Duration, error) {
	valueOf := reflect.ValueOf(t)
	value := valueOf.FieldByName(key).Interface().(int64)
	return time.Second * time.Duration(value), nil
}

package config

import (
	"reflect"
	"register-go/infra/utils/common"
)

// 日志配置
type LogConfig struct {
	Level        string `yaml:"level"`        // 日志级别
	LogFileName  string `yaml:"logFileName"`  // 日志文件名称
	FilePath     string `yaml:"filePath"`     // 日志文件路径
	MaxAge       int    `yaml:"maxAge"`       // 日志保存时间
	RotationTime int    `yaml:"rotationTime"` // 日志切割时间间隔
}

func (l LogConfig) GetString(attr string) string {
	valueOf := reflect.ValueOf(l)
	value := valueOf.FieldByName(attr).Interface().(string)
	if common.StrIsBlank(value) {
		return common.NilString
	}
	return value
}

func (l LogConfig) GetStringDefault(attr, defaultValue string) string {
	valueOf := reflect.ValueOf(l)
	value := valueOf.FieldByName(attr).Interface().(string)
	if common.StrIsBlank(value) {
		return defaultValue
	}
	return value
}

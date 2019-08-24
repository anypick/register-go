package config

import (
	"reflect"
	"register-go/infra/utils/common"
	"time"
)

type MySqlConfig struct {
	DriverName      string        `yaml:"driverName"`      // 驱动名称
	IpAddr          string        `yaml:"ipAddr"`          // ip地址
	Port            string        `yaml:"port"`            // 端口
	Username        string        `yaml:"username"`        // 用户名
	Password        string        `yaml:"password"`        // 密码
	Database        string        `yaml:"database"`        // 数据库名称
	MaxOpenConn     int           `yaml:"maxOpenConn"`     // 最大连接数
	MaxIdeConn      int           `yaml:"maxIdeConn"`      // 最大等待连接
	ConnMaxLifetime time.Duration `yaml:"connMaxLifetime"` // 连接最大存活时间
}

func (m MySqlConfig) GetStringByDefault(fieldName, defaultValue string) string {
	stringValue := reflect.ValueOf(m).FieldByName(fieldName).Interface().(string)
	if common.StrIsBlank(stringValue) {
		return defaultValue
	}
	return stringValue
}

func (m MySqlConfig) GetIntByDefault(fieldName string, defaultValue int) int {
	intValue := reflect.ValueOf(m).FieldByName(fieldName).Interface().(int)
	if intValue == 0 {
		return defaultValue
	}
	return intValue
}

func (m MySqlConfig) GetDurationDefault(fieldName string, defaultValue time.Duration) time.Duration {
	durationValue := reflect.ValueOf(m).FieldByName(fieldName).Interface().(time.Duration)
	if durationValue == 0 {
		return time.Second * defaultValue
	}
	return time.Second * durationValue
}

package props

import "time"

type ConfigSource interface {
	Strings(key string) []string
	Ints(key string) []int
	Float64s(key string) []float64
	Durations(key string) []time.Duration
	//
	Get(key string) (string, error)
	GetDefault(key, defaultValue string) string

	//
	GetInt(key string) (int, error)
	GetIntDefault(key string, defaultValue int) int
	//
	GetDuration(key string) (time.Duration, error)
	GetDurationDefault(key string, defaultValue time.Duration) time.Duration
	//
	GetBool(key string) (bool, error)
	GetBoolDefault(key string, defaultValue bool) bool
	//
	GetFloat64(key string) (float64, error)
	GetFloat64Default(key string, defaultValue float64) float64

	//t必须为指针型
	Unmarshal(t interface{}) error
}

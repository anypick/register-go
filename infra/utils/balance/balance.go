package balance

type Balance interface {
	Next(key string) int
}

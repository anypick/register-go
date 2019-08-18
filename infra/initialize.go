package infra

// 用于业务的代码的注入，例如Dao层，Service层，Controller层
type Initializer interface {
	Init()
}

type InitializerRegister struct {
	Initializers []Initializer
}

func (i *InitializerRegister) Register(ai Initializer) {
	i.Initializers = append(i.Initializers, ai)
}

var apiRegister *InitializerRegister = new(InitializerRegister)

func RegisterApi(ai Initializer) {
	apiRegister.Register(ai)
}

func GetApiRegister() []Initializer {
	return apiRegister.Initializers
}

// 用于注册业务结构体
type BaseInitializerStarter struct {
	BaseStarter
}

func (b *BaseInitializerStarter) Setup(ctx StarterContext) {
	for _, register := range GetApiRegister() {
		register.Init()
	}
}

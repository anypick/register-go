package basegin

import "github.com/anypick/register-go/infra"


var apiRegister *infra.InitializerRegister = new(infra.InitializerRegister)

// 注册业务类
func RegisterApi(ai infra.Initializer) {
	apiRegister.Register(ai)
}

// 获取所有注册业务类
func GetApiRegister() []infra.Initializer {
	return apiRegister.Initializers
}

type WebStarter struct {
	infra.BaseStarter
}

func (w *WebStarter) Setup(ctx infra.StarterContext) {
	for _, register := range GetApiRegister() {
		register.Init()
	}
}

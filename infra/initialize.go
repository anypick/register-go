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


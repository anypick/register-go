package infra

import (
	"github.com/sirupsen/logrus"
	"reflect"
	"register-go/infra/utils/props"
)

// 负责Starter各个阶段方法的调用
type BootApplication struct {
	conf           props.YamlSource
	starterContext StarterContext
}

func New(conf props.YamlSource) *BootApplication {
	application := &BootApplication{conf, StarterContext{}}
	application.starterContext[defaultProps] = conf
	return application
}

func (b *BootApplication) Start() {
	//1. 初始化starter
	b.init()
	//2. 安装starter
	b.setup()
	//3. 启动starter
	b.start()
}

func (b *BootApplication) init() {
	logrus.Info("Application init...")
	starters := StarterRegister.AllStarters()
	for _, starter := range starters {
		starter.Init(b.starterContext)
	}
}

func (b *BootApplication) setup() {
	logrus.Info("Application setup...")
	starters := StarterRegister.AllStarters()
	for _, starter := range starters {
		starter.Setup(b.starterContext)
	}
}

func (b *BootApplication) start() {
	logrus.Info("Application start")
	starters := StarterRegister.AllStarters()
	for index, starter := range starters {
		typ := reflect.TypeOf(starter)
		logrus.Info("Starting: ", typ.String())
		if starter.StartBlocking() {
			if index+1 == len(StarterRegister.AllStarters()) {
				starter.Start(b.starterContext)
			} else {
				go starter.Start(b.starterContext)
			}
		} else {
			starter.Start(b.starterContext)
		}
	}
}

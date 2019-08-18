package infra

import (
	"register-go/infra/utils/props"
	"sort"
)

// 启动优先级
type PriorityGroup int

const (
	// 配置文件名称, TODO 支持环境隔离
	configFilepathName = "resources/application.yml"
	// 配置文件key，通过这key，可以获取配置文件对象
	defaultProps                      = "props"
	BasicResourcesGroup PriorityGroup = 20
	INT_MAX                           = int(^uint(0) >> 1)
	DEFAULT_PRIORITY                  = 0
)

// 加载配置文件
type StarterContext map[string]interface{}

// 获取yaml配置文件
func (s StarterContext) Yaml() props.YamlSource {
	yml := s[defaultProps]
	return yml.(props.YamlSource)
}

type Starter interface {
	// 初始化一些资源,如配置加载
	Init(StarterContext)

	// 系统基础资源的安装，让资源达到可用的状态，但是还没有被使用
	Setup(StarterContext)

	// 启动基础资源，资源正在使用，例如开启web服务器，
	Start(StarterContext)

	// 启动器是否可阻塞， 这里阻塞的需要放到最后
	StartBlocking() bool

	// 资源停止和销毁：
	// 通常在启动时遇到异常时或者启用远程管理时，用于释放资源和终止资源的使用，
	// 通常要优雅的释放，等待正在进行的任务继续，但不再接受新的任务
	Stop(StarterContext)

	// 定义启动的优先级，越大优先级越高
	PriorityGroup() PriorityGroup
	Priority() int
}

// 基础Starter, 通过组合，子结构体不需要全部实现父接口的方法
type BaseStarter struct{}

func (b *BaseStarter) Init(StarterContext)          {}
func (b *BaseStarter) Setup(StarterContext)         {}
func (b *BaseStarter) Start(StarterContext)         {}
func (b *BaseStarter) StartBlocking() bool          { return false }
func (b *BaseStarter) Stop(StarterContext)          {}
func (b *BaseStarter) PriorityGroup() PriorityGroup { return BasicResourcesGroup }
func (b *BaseStarter) Priority() int                { return DEFAULT_PRIORITY }

// 集中管理Starter
type starterRegister struct {
	blockStarter     []Starter
	noneBlockStarter []Starter
}

func (s *starterRegister) Register(starter Starter) {
	if starter.StartBlocking() {
		s.blockStarter = append(s.blockStarter, starter)
	} else {
		s.noneBlockStarter = append(s.noneBlockStarter, starter)
	}
}

func (s *starterRegister) AllStarters() []Starter {
	starters := make([]Starter, 0)
	// 非阻塞在前面
	starters = append(starters, s.noneBlockStarter...)
	// 阻塞的在后面
	starters = append(starters, s.blockStarter...)
	// 按照优先级进行排序，排序只对Start时候有效
	sortStarter(starters)
	return starters
}

var StarterRegister = new(starterRegister)

// 由外部调用，将starter注册到starterRegister中，由容器统一管理，starter的实例
func Register(s Starter) {
	StarterRegister.Register(s)
}

// starter实现排序, 需要实现sort.Interface接口
type Starters []Starter

func (s Starters) Len() int {
	return len(s)
}

func (s Starters) Less(i, j int) bool {
	return s[i].PriorityGroup() > s[j].PriorityGroup() && s[i].Priority() > s[j].Priority()
}

func (s Starters) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func sortStarter(starters []Starter) {
	sort.Sort(Starters(starters))
}

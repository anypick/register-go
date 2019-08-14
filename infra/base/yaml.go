package base

import (
	"github.com/anypick/register-go/infra"
	"github.com/anypick/register-go/infra/utils/props"
)

var (
	yamlProps props.YamlSource
)

func YamlProps() props.YamlSource {
	return yamlProps
}

type YamlStarter struct {
	infra.BaseStarter
}

func (p *YamlStarter) Init(ctx infra.StarterContext) {
	yamlProps = ctx.Yaml()
}



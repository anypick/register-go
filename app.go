package register_go

import (
	_ "github.com/anypick/register-go/api/example"
	"github.com/anypick/register-go/infra"
	"github.com/anypick/register-go/infra/base"
	"github.com/anypick/register-go/infra/base/gin"
)

func init() {
	infra.Register(&base.YamlStarter{})
	infra.Register(&basegin.GinStarter{})
	infra.Register(&basegin.WebStarter{})
}

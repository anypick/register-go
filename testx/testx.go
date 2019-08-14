package testx

import (
	"fmt"
	"github.com/anypick/register-go/infra"
	"github.com/anypick/register-go/infra/base"
	"github.com/anypick/register-go/infra/base/redis"
	"github.com/anypick/register-go/infra/utils/props"
)

func init() {
	infra.Register(&base.YamlStarter{})
	infra.Register(&baseredis.RedisReplicationStarter{})

	banner := `
                     .__         __                                        
_______  ____   ____ |__| ______/  |_  ___________            ____   ____  
\_  __ \/ __ \ / ___\|  |/  ___|   __\/ __ \_  __ \  ______  / ___\ /  _ \ 
 |  | \|  ___// /_/  >  |\___ \ |  | \  ___/|  | \/ /_____/ / /_/  >  <_> )
 |__|   \___  >___  /|__/____  >|__|  \___  >__|            \___  / \____/ 
            \/_____/         \/           \/               /_____/
`
	fmt.Println(banner)
	yamlConf := props.NewYamlSource("../../resources/application.yml")
	infra.New(*yamlConf).Start()
}

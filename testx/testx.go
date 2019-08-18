package testx

import (
	"fmt"
	"register-go/infra"
	"register-go/infra/base"
	"register-go/infra/base/redis"
	"register-go/infra/utils/props"
)

func init() {
	infra.Register(&base.YamlStarter{})
	infra.Register(&baseredis.RedisReplicationStarter{})
	infra.Register(&infra.BaseInitializerStarter{})
	banner := `
                     .__         __                                        
_______  ____   ____ |__| ______/  |_  ___________            ____   ____  
\_  __ \/ __ \ / ___\|  |/  ___|   __\/ __ \_  __ \  ______  / ___\ /  _ \ 
 |  | \|  ___// /_/  >  |\___ \ |  | \  ___/|  | \/ /_____/ / /_/  >  <_> )
 |__|   \___  >___  /|__/____  >|__|  \___  >__|            \___  / \____/ 
            \/_____/         \/           \/               /_____/
`
	fmt.Println(banner)
	//yamlConf := props.NewYamlSource("../../resources/application.yml")
	yamlConf := props.NewYamlSource("resources/application.yml")
	infra.New(*yamlConf).Start()
}

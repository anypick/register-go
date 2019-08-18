package main

import (
	"fmt"
	_ "register-go"
	"register-go/infra"
	"register-go/infra/utils/props"
	_ "register-go/src/site/controller"
)

func main() {
	banner := `
                     .__         __                                        
_______  ____   ____ |__| ______/  |_  ___________            ____   ____  
\_  __ \/ __ \ / ___\|  |/  ___|   __\/ __ \_  __ \  ______  / ___\ /  _ \ 
 |  | \|  ___// /_/  >  |\___ \ |  | \  ___/|  | \/ /_____/ / /_/  >  <_> )
 |__|   \___  >___  /|__/____  >|__|  \___  >__|            \___  / \____/ 
            \/_____/         \/           \/               /_____/
`
	fmt.Println(banner)
	yamlConf := props.NewYamlSource("resources/application.yml")
	application := infra.New(*yamlConf)
	application.Start()
}

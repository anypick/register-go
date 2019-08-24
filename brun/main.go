package main

import (
	"flag"
	"fmt"
	_ "register-go"
	"register-go/infra"
	"register-go/infra/utils/common"
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
	profile := flag.String("profile", "", "环境信息")
	flag.Parse()
	resource := ""
	if common.StrIsBlank(*profile) {
		resource = "resources/application.yml"
	} else {
		resource = fmt.Sprintf("resources/application-%s.yml", *profile)
	}
	fmt.Println(resource)
	yamlConf := props.NewYamlSource(resource)
	application := infra.New(*yamlConf)
	application.Start()
}

package example

import (
	"github.com/anypick/register-go/infra/base/gin"
	"github.com/gin-gonic/gin"
)

func init() {
	basegin.RegisterApi(new(TestController))
}

type TestController struct {
}

func (t *TestController) Init() {
	gin := basegin.Gin()
	group := gin.Group("/v1/test")
	group.GET("/hello", t.hello)
}


func (t *TestController) hello (ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "pong",
	})
}

package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"register-go/infra"
	"register-go/infra/base/gin"
	"register-go/src/site/service"
)

func init() {
	infra.RegisterApi(&SiteController{})
}

type SiteController struct {
	siteService service.ISiteService
}

func (c *SiteController) Init() {
	c.siteService = service.GetSiteService()
	app := basegin.Gin().Group("/v1/site")
	app.POST("/add", c.Add)
	app.POST("/get", c.Get)
}

func (c *SiteController) Add(ctx *gin.Context) {
	flag := c.siteService.Add()
	if flag.Success {
		ctx.String(http.StatusOK, "成功")
	} else {
		ctx.String(http.StatusOK, "失败")
	}
}

func (c *SiteController) Get(ctx *gin.Context) {
	data := c.siteService.Get()
	ctx.JSON(200, data)
}

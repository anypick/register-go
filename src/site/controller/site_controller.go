package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"register-go/infra"
	"register-go/infra/base/gin"
	"register-go/infra/redisutil"
	"register-go/infra/utils/common"
	"register-go/src/site/dto"
	"register-go/src/site/service"
	"strconv"
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
	app.POST("/get/:siteId/:langCode", c.GetById)
	app.POST("/getByField", c.GetByField)
}

/**
json请求：
{
   "Id":1001,
   "SiteId": 10001,
   "SiteCode":"S1001",
   "SiteId":10001,
   "SiteName":"上海",
}
 */
func (c *SiteController) Add(ctx *gin.Context) {
	// 获取参数
	var site dto.SiteDto
	if err := ctx.Bind(&site); err == nil {
		resp := c.siteService.Add(site)
		ctx.JSON(http.StatusOK, resp)
	} else {
		ctx.JSON(http.StatusInternalServerError, common.NewRespFailWithMsg("参数格式错误"))
	}
}

// 请求参数 /v1/site/get/10001/zh_CN
func (c *SiteController) GetById(ctx *gin.Context) {
	var (
		siteId int64
		err    error
	)
	langCode := ctx.Param("langCode")
	if siteId, err = strconv.ParseInt(ctx.Param("siteId"), 10, 64); err != nil || siteId == 0 {
		ctx.JSON(http.StatusOK, common.NewRespFailWithMsg("参数类型错误或者参数不存在"))
	} else if common.StrIsBlank(langCode) {
		langCode = redisutil.DefaultLangCode
	}
	ctx.JSON(http.StatusOK, c.siteService.GetById(siteId, langCode))
}

// 表单请求：page: 1, pageSize: 10, fieldName: siteName, fieldValue: 10001
func (c *SiteController) GetByField(ctx *gin.Context) {
	var (
		page, ePage         = strconv.Atoi(ctx.DefaultPostForm("page", "1"))
		pageSize, ePageSize = strconv.Atoi(ctx.DefaultPostForm("pageSize", "10"))
		fieldName           = ctx.PostForm("fieldName")
		fieldValue          = ctx.PostForm("fieldValue")

	)
	if ePage != nil || ePageSize != nil {
		logrus.Error(ePage, ePageSize)
	}

	if common.StrIsBlank(fieldName) {
		ctx.JSON(http.StatusOK, common.NewRespFailWithMsg("参数错误"))
		return
	}
	if page == 0 {
		page = redisutil.DefaultPage
	}
	if pageSize == 0 {
		pageSize = redisutil.DefaultPageSize
	}
	ctx.JSON(http.StatusOK, c.siteService.GetByField(fieldName, fieldValue, page, pageSize))
}

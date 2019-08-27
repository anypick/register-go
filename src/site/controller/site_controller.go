package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"register-go/infra"
	"register-go/infra/base/gin"
	"register-go/infra/base/redis"
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
	siteService    service.ISiteService
	siteSqlService *service.SiteSqlService
}

func (c *SiteController) Init() {
	c.siteService = service.GetSiteService()
	c.siteSqlService = service.GetSiteSqlService()
	app := basegin.Gin().Group("/v1/site")
	app.POST("/redis/add", c.Add)
	app.POST("/redis/get/:siteId/:langCode", c.GetById)
	app.POST("/redis/getByField", c.GetByField)
	app.POST("/sql/insert", c.InsertSql)
	// 插入hashes数据类型到redis
	app.POST("/redis/hash/add", c.AddHash)
	app.GET("/redis/hash/get/:siteId/:langCode", c.GetHashById)
	app.GET("/redis/client/get", c.TestRedisClient)
	app.GET("/redis/cluster/get", c.TestRedisClusterClient)
	app.GET("/redis/hash/getall/:langCode", c.GetAllHash)
	app.POST("/rabbit/direct/sendMsg", c.SendMsg)
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

/**
[{
  "Id": 10004,
  "SiteId": 10005,
  "SiteCode":"S1005",
  "SiteName":"人民广场店xx"
},{
  "Id": 10003,
  "SiteId": 10003,
  "SiteCode":"S1003",
  "SiteName":"浦东店铺"
},{
  "id": 10002,
  "SiteId": 10006,
  "SiteCode":"S1006",
  "SiteName":"人民广场店xxx"
}]
*/
func (c *SiteController) InsertSql(ctx *gin.Context) {
	var sites []dto.SiteDto
	if err := ctx.BindJSON(&sites); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusOK, common.NewRespFail())
		return
	}
	ctx.JSON(http.StatusOK, c.siteSqlService.UpdateOrInsert(sites))
}

/**
{
   "Id":1003,
   "SiteCode":"S1003",
   "SiteId":10003,
   "SiteName":"北京店"
}
*/
func (c *SiteController) AddHash(ctx *gin.Context) {
	var site dto.SiteDto
	if err := ctx.BindJSON(&site); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusOK, common.NewRespFail())
		return
	}
	ctx.JSON(http.StatusOK, c.siteService.AddHash(site))
}

func (c *SiteController) GetHashById(ctx *gin.Context) {
	var (
		langCode    = ctx.Param("langCode")
		siteIdParam = ctx.Param("siteId")
		siteId      int64
		err         error
	)
	if common.StrIsBlank(siteIdParam) {
		ctx.JSON(http.StatusOK, common.NewRespFailWithMsg("参数错误"))
		return
	}
	if siteId, err = strconv.ParseInt(siteIdParam, 10, 64); err != nil {
		ctx.JSON(http.StatusOK, common.NewRespFailWithMsg("参数解析错误"))
		return
	}
	ctx.JSON(http.StatusOK, c.siteService.GetHashById(siteId, langCode))
}

// 表单参数：page:1， pageSize:10
// url参数：/:langCode
func (c *SiteController) GetAllHash(ctx *gin.Context) {
	var (
		page, ePage         = strconv.Atoi(ctx.DefaultQuery("page", "1"))
		pageSize, ePagesize = strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
		langCode            = ctx.Param("langCode")
	)
	if ePage != nil || ePagesize != nil {
		logrus.Error(ePage, ePagesize)
		ctx.JSON(http.StatusOK, common.NewRespFailWithMsg("参数错误"))
		return
	}
	ctx.JSON(http.StatusOK, c.siteService.GetAllHash(page, pageSize, langCode))
}

// 发送RabbitMQ消息
func (c *SiteController) SendMsg(ctx *gin.Context) {
	var site dto.SiteDto

	if err := ctx.Bind(&site); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusOK, common.NewRespFailWithMsg("参数错误"))
		return
	}
	ctx.JSON(http.StatusOK, c.siteSqlService.SendMsg(site))
}

// 用于测试Redis客户端
func (c *SiteController) TestRedisClient(ctx *gin.Context) {
	client := baseredis.RedisClient(baseredis.Sentinel)
	fmt.Println(client.Ping())
	fmt.Println(client.HGet("hmall:Site:zh_CN:site", "10005").Val())
	ctx.String(http.StatusOK, "hao")
}


// 用于测试Redis客户端
func (c *SiteController) TestRedisClusterClient(ctx *gin.Context) {
	client := baseredis.GetRedisCluster()
	fmt.Println(client.Ping())
	ctx.String(http.StatusOK, "hao")
}
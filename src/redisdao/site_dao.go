package redisdao

import (
	"fmt"
	"register-go/infra"
	"register-go/infra/redisutil"
)

func init() {
	infra.RegisterApi(&SiteDao{})
}

var siteDao *SiteDao

func GetSiteDao() *SiteDao {
	return siteDao
}

type SiteDao struct {
	redisutil.BaseDao
}

func (s *SiteDao) Init() {
	siteDao = &SiteDao{}
	siteDao.Catalog = "Site"
	siteDao.Clazz = "site"
	siteDao.IdDesc = redisutil.FieldDescriptor{FieldName: "SiteId", FieldType: redisutil.TypeEq}
	siteDao.CreateFieldDescriptor("SiteName", redisutil.TypeMatch)
	siteDao.CreateFieldDescriptor("SiteCode", redisutil.TypeEq)
	// 开业时间进行排序
	siteDao.CreateFieldDescriptor("CreateTime", redisutil.TypeRange)
}

func (s *SiteDao) InitTest() {
	fmt.Println("就是测试")
}

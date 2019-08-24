package service

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"register-go/infra"
	"register-go/infra/redisutil"
	"register-go/infra/utils/common"
	"register-go/src/redisdao"
	"register-go/src/site/dto"
	"time"
)

var siteService ISiteService

func GetSiteService() ISiteService {
	return siteService
}

func init() {
	infra.RegisterApi(&SiteServiceImpl{})
}

type ISiteService interface {
	Add(site dto.SiteDto) common.ResponseData
	GetById(siteId int64, langCode string) common.ResponseData
	GetByField(fieldName string, fieldValue interface{}, page int, pageSize int) common.ResponseData
	AddHash(site dto.SiteDto) common.ResponseData
	GetHashById(siteId int64, langCode string) common.ResponseData
	GetAllHash(page int, pageSize int, langCode string) common.ResponseData
}

type SiteServiceImpl struct {
	siteDao redisutil.IBaseDao
}

func (s *SiteServiceImpl) Init() {
	siteService = &SiteServiceImpl{siteDao: redisdao.GetSiteDao()}
}

func (s *SiteServiceImpl) Add(site dto.SiteDto) common.ResponseData {
	site.UpdateTime = time.Now()
	site.CreateTime = time.Now()
	b, e := s.siteDao.Add(common.StructToMap(site), 0, redisutil.DefaultLangCode)
	if b && e == nil {
		resp := common.NewRespSucc()
		return resp
	} else {
		resp := common.NewRespFailWithMsg("添加数据错误")
		return resp
	}
}

func (s *SiteServiceImpl) GetById(siteId int64, langCode string) common.ResponseData {
	data, err := s.siteDao.Get(siteId, langCode)
	if err != nil || data == nil {
		logrus.Error(err)
		return common.NewRespSucc()
	}
	respData := []interface{}{data}
	return common.NewRespSuccWithData(respData, 1)
}

func (s *SiteServiceImpl) GetByField(fieldName string, fieldValue interface{}, page int, pageSize int) common.ResponseData {
	var (
		data []map[string]interface{}
		err  error
	)
	if data, err = s.siteDao.GetByField(fieldValue, fieldName, "", page, pageSize); err != nil {
		logrus.Info(err)
		return common.NewRespSucc()
	}
	return common.NewRespSuccWithData(data, len(data))
}
func (s *SiteServiceImpl) AddHash(site dto.SiteDto) common.ResponseData {
	site.CreateTime = time.Now()
	site.UpdateTime = time.Now()
	if flag, err := s.siteDao.AddHash(common.StructToMap(site), 0, common.NilString); err != nil || !flag {
		logrus.Error("插入数据失败", err)
		return common.NewRespFail()
	}
	return common.NewRespSucc()
}

func (s *SiteServiceImpl) GetHashById(siteId int64, langCode string) common.ResponseData {
	var (
		dataStr string
		data    dto.SiteDto
		err     error
	)
	if dataStr, err = s.siteDao.GetHash(siteId, langCode); err != nil {
		logrus.Error(err)
		common.NewRespSucc()
	}
	if err = json.Unmarshal([]byte(dataStr), &data); err != nil {
		logrus.Error(err)
		common.NewRespFailWithMsg("内部错误")
	}
	return common.NewRespSuccWithData([]dto.SiteDto{data}, 1)
}

func (s *SiteServiceImpl) GetAllHash(page int, pageSize int, langCode string) common.ResponseData {
	allData, err := s.siteDao.GetAllHash(page, pageSize, langCode)
	if err != nil {
		logrus.Error("err")
		return common.NewRespFail()
	}
	return common.NewRespSuccWithData(allData, len(allData))
}

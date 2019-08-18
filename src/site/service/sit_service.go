package service

import (
	"register-go/infra"
	"register-go/infra/redisutil"
	"register-go/infra/utils/common"
	"register-go/src/redisdao"
)

var siteService ISiteService

func GetSiteService() ISiteService {
	return siteService
}

func init() {
	infra.RegisterApi(&SiteServiceImpl{})
}

type ISiteService interface {
	Add() common.ResponseData
	Get() common.ResponseData
}

type SiteServiceImpl struct {
	siteDao redisutil.IBaseDao
}

func (s *SiteServiceImpl) Init() {
	siteService = &SiteServiceImpl{siteDao: redisdao.GetSiteDao()}
}

func (s *SiteServiceImpl) Add() (resp common.ResponseData) {
	resp = common.NewRespSucc()
	return
}

func (s *SiteServiceImpl) Get() common.ResponseData {

	data, err := s.siteDao.Get(10008, redisutil.DefaultLangCode)
	if err != nil || data == nil {
		return common.NewRespSucc()
	}
	respData := make([]interface{}, 1)
	respData[0] = data
	return common.NewRespSuccWithData(respData)
}

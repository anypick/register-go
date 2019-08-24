package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"register-go/infra"
	"register-go/infra/base/mysql"
	"register-go/infra/utils/common"
	"register-go/src/site/dao"
	"register-go/src/site/dto"
)

var siteSqlService *SiteSqlService

func GetSiteSqlService() *SiteSqlService {
	return siteSqlService
}

func init() {
	infra.RegisterApi(&SiteSqlService{})
}

var siteSqlDao = dao.GetSitSqlDao()

type SiteSqlService struct {
}

func (s *SiteSqlService) Init() {
	siteSqlService = s
}

func (s *SiteSqlService) UpdateOrInsert(sites []dto.SiteDto) common.ResponseData {
	var (
		err error
	)
	err = basesql.DbTxRunner(func(runner *basesql.Runner) error {
		ctx := basesql.WithValueContext(context.TODO(), runner)
		for _, site := range sites {
			// update
			if site.Id != 0 {
				if err = siteSqlDao.UpdateById(site, ctx); err != nil {
					logrus.Error(err, site)
					return err
				}
			} else { //插入
				if err = siteSqlDao.InsertOne(site, ctx); err != nil {
					logrus.Error(err)
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return common.NewRespFailWithMsg("更新或插入失败")
	}
	return common.NewRespSucc()
}
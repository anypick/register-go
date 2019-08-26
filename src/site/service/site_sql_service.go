package service

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"register-go/infra"
	"register-go/infra/base/mysql"
	"register-go/infra/utils/common"
	"register-go/src/rabbitmq"
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
var siteRabbit = rabbitmq.GetSiteRabbit()

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

func (s *SiteSqlService) SendMsg(site dto.SiteDto) common.ResponseData {
	var (
		siteData []byte
		err error
	)
	siteRabbit = rabbitmq.GetSiteRabbit()
	if siteData, err = json.Marshal(site); err != nil {
		logrus.Error(err)
		return common.NewRespFail()
	}
	// 过期消息
	if err := siteRabbit.Send(amqp.Publishing{Body:siteData, Expiration: "10000"}); err != nil {
		logrus.Error(err)
		return common.NewRespFail()
	}
	return common.NewRespSucc()
}


package impl

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"register-go/infra/base/grpc"
	"register-go/infra/utils/common"
	"register-go/src/grpc"
	"register-go/src/site/dto"
)

// ====================================================================================
// ==============================充当服务端代码===================================
// ====================================================================================

func init() {
	basegrpc.RegisterGrpc(&SiteService{})
}

type SiteService struct{}

func (s *SiteService) RegisterGrpc(rpc *grpc.Server) {
	site_service.RegisterSiteServiceServer(rpc, s)
}

func (s *SiteService) GetSiteDetailBySiteCode(ctx context.Context, param *site_service.SiteDto) (*site_service.ResponseData, error) {
	fmt.Println(param)
	return &site_service.ResponseData{Code: 1001, Msg: "访问成功"}, nil
}

// ====================================================================================
// ==============================访问其他服务端代码===================================
// ====================================================================================

func TestGrpc(dto dto.SiteDto) common.ResponseData {
	var conn = basegrpc.GetGrpcConn("site_service")
	var (
		siteDto = site_service.SiteDto{SiteCode: dto.SiteCode}
		client  = site_service.NewSiteServiceClient(conn)
		data    *site_service.ResponseData
		err     error
	)
	if data, err = client.GetSiteDetailBySiteCode(context.TODO(), &siteDto); err != nil {
		logrus.Error(err)
		return common.NewRespFail()
	}
	fmt.Println(data)
	return common.NewRespSucc()
}

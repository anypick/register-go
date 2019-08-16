package redisdao

import (
	"github.com/anypick/register-go/infra/redisutil"
)

func init() {
	//infra.RegisterApi()
}


type SiteDao struct {
	redisutil.BaseDao
}

func (s *SiteDao) Init() {

}

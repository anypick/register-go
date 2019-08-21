package dto

import (
	"time"
)

type SiteDto struct {
	Id         int64     `form:"Id"`
	SiteId     int64     `form:"SiteId"`
	SiteCode   string    `form:"SiteCode"`
	SiteName   string    `form:"SiteName"`
	CreateTime time.Time `form:"CreateTime" time_format:"2006-01-02" time_utc:"1"`
	UpdateTime time.Time `form:"UpdateTime" time_format:"2006-01-02" time_utc:"1"`
}

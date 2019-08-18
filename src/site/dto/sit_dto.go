package dto

import "time"

type SiteDto struct {
	Id         int64
	SiteId     int64
	SiteCode   string
	SiteName   string
	CreateTime time.Time
	UpdateTime time.Time
}

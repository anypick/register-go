package dao

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"register-go/infra"
	"register-go/infra/base/mysql"
	"register-go/src/site/dto"
)

var siteDao *SitSqleDao

func init() {
	infra.RegisterApi(&SitSqleDao{})
}

func GetSitSqlDao() *SitSqleDao {
	return siteDao
}

type SitSqleDao struct{}

func (s *SitSqleDao) Init() {
	siteDao = s
}

func (s *SitSqleDao) InsertOne(site dto.SiteDto, ctx context.Context) error {
	var (
		err       error
		insertSql = `INSERT INTO SITE(site_id, site_code, site_name) VALUES (?, ?, ?)`
		result    sql.Result
		lastId    int64
	)
	err = basesql.ExecuteContext(ctx, func(runner *basesql.Runner) error {
		if result, err = runner.Tx.Exec(insertSql, site.SiteId, site.SiteCode, site.SiteName); err != nil {
			return err
		}
		if lastId, err = result.LastInsertId(); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *SitSqleDao) UpdateById(site dto.SiteDto, ctx context.Context) error {
	var (
		err       error
		updateSql = `UPDATE site SET site_name = ?, site_code = ? WHERE id = ? and site_id = ?`
		result    sql.Result
	)
	err = basesql.ExecuteContext(ctx, func(runner *basesql.Runner) error {
		if result, err = runner.Tx.Exec(updateSql, site.SiteName, site.SiteCode, site.Id, site.SiteId); err != nil {
			return err
		}
		if count, err := result.RowsAffected(); err != nil {
			return err
		} else if count == 1 {
			return nil
		}
		return errors.New("更新失败")
	})
	return err
}

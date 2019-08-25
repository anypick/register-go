package basesql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"register-go/infra"
	"time"
)

var db *sql.DB

func GetDb() *sql.DB {
	return db
}

type MySqlStarter struct {
	infra.BaseStarter
}

func (d *MySqlStarter) Init(ctx infra.StarterContext) {
	config := ctx.Yaml().MySqlConfig
	var err error
	// user:password@tcp(127.0.0.1:3306)/test
	db, err = sql.Open(config.DriverName, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", config.Username, config.Password, config.IpAddr, config.Port, config.Database))
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(config.GetDurationDefault("ConnMaxLifetime", time.Duration(7*24*60*60)))
	db.SetMaxOpenConns(config.GetIntByDefault("MaxOpenConn", 1000))
	db.SetMaxIdleConns(config.GetIntByDefault("MaxIdeConn", 1000))
	if db.Ping() == nil {
		logrus.Info("mysql up")
	} else {
		logrus.Warn(db.Ping())
	}
}

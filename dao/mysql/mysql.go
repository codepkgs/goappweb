package mysql

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/x-hezhang/gowebapp/settings"
)

var db *sqlx.DB

func Init() (err error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		settings.Conf.MySQLConfig.User,
		settings.Conf.MySQLConfig.Password,
		settings.Conf.MySQLConfig.Host,
		settings.Conf.MySQLConfig.Port,
		settings.Conf.MySQLConfig.Database,
	)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return
	}
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute)
	return
}

func Close() {
	_ = db.Close()
}

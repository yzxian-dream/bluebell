package mysql

import (
	"fmt"
	"webapp/settings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

// Init 初始化MySQL连接
func InitDB(cfg *settings.MysqlConf) (err error) {
	// "user:password@tcp(host:port)/dbname"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBname)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return
	}
	db.SetMaxOpenConns(cfg.Max_Conns)
	db.SetMaxIdleConns(cfg.Max_Idle_Conns)
	return
}

// Close 关闭MySQL连接
func Close() {
	_ = db.Close()
}

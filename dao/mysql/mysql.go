package mysql

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"strings"
	"time"
	"webapp/settings"
)

var db *sqlx.DB

// Init 初始化MySQL连接
func InitDB(cfg *settings.MysqlConf) (err error) {
	// "user:password@tcp(host:port)/dbname"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBname)
	//db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return
	}
	db, err = sqlx.Connect("mysql", strings.Split(dsn, "/")[0]+"/")
	if err != nil {
		return errors.Wrap(err, "fail to open mysql")
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+"bluebell")
	if err != nil {
		fmt.Printf("Error %s when creating DB\n", err)
		return
	}
	no, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("Error %s when fetching rows", err)
		return
	}
	db, err = sqlx.Open("mysql", dsn)
	if err != nil {
		return errors.Wrap(err, "fail to open mysql")
	}
	err = db.Ping()
	if err != nil {
		return errors.Wrap(err, "fail to connect to mysql")
	}
	fmt.Printf("rows affected %d\n", no)
	db.SetMaxOpenConns(cfg.Max_Conns)
	db.SetMaxIdleConns(cfg.Max_Idle_Conns)
	return
}

// Close 关闭MySQL连接
func Close() {
	_ = db.Close()
}

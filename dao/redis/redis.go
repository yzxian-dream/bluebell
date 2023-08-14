package redis

import (
	"fmt"
	"webapp/settings"

	"github.com/go-redis/redis"
)

var rdb *redis.Client

func Init(cfg *settings.RedisConf) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,  // 密码
		DB:       cfg.DB,        // 数据库
		PoolSize: cfg.Pool_Size, // 连接池大小
	})
	_, err = rdb.Ping().Result()
	return

}

func Close() {
	_ = rdb.Close()
}

package redis

import (
	"context"
	"os"
	"log"

	rv8 "github.com/go-redis/redis/v8"
)

var (
//	RedisCtx = context.Background()
)

func NewRedis(ctx context.Context, host string, port string, password string, database int) *rv8.Client {
	// 异常处理
	defer func() {
		if err := recover(); err != nil {
			log.Println("Redis Err:", err)
			os.Exit(1)
		}
	}()

	var conn *rv8.Client
	conn = rv8.NewClient(&rv8.Options{
		Addr:     host + ":" + port,
		Password: password, // no password set
		DB:       database, // use default DB
	})
	_, err := conn.Ping(ctx).Result() // 检查Redis链接
	if err != nil {
		panic(err)
	}

	return conn
}

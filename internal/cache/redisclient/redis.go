package redisclient

import (
	"app/internal/config"
	"crypto/tls"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func Load() {
	redisCfg := config.GetConfig().Redis
	RedisClient = NewRedisClient(redisCfg.Endpoint, redisCfg.User, redisCfg.Password, redisCfg.TLS)
}

func NewRedisClient(Addr string, user string, password string, TLS bool) *redis.Client {

	if TLS {
		redisClient := redis.NewClient(&redis.Options{
			Addr:      Addr,     // AWS Redis 地址或本地地址
			Password:  password, // 如果没有密码，设为 ""
			Username:  user,
			DB:        0,             // Redis 默认 DB 0
			TLSConfig: &tls.Config{}, // 启用 TLS，忽略证书验证
		})
		// 启用 tracing
		if err := redisotel.InstrumentTracing(redisClient); err != nil {
			panic(err)
		}
		return redisClient
	} else {
		redisClient := redis.NewClient(&redis.Options{
			Addr:     Addr,     // AWS Redis 地址或本地地址
			Password: password, // 如果没有密码，设为 ""
			Username: user,
			DB:       0, // Redis 默认 DB 0
		})
		// 启用 tracing
		if err := redisotel.InstrumentTracing(redisClient); err != nil {
			panic(err)
		}
		return redisClient
	}
}

package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var RDB *redis.Client

func Init() (err error) {
	RDB = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			viper.GetString("redis.host"),
			viper.GetInt("redis.port"),
		),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
		PoolSize: viper.GetInt("redis.pool_size"),
	})
	_, err = RDB.Ping().Result()
	return err
}

func Close() {
	_ = RDB.Close()
}

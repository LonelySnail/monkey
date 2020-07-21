package util

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

func GetRedisPool(url string) (*redis.Pool, error) {
	pool := &redis.Pool{
		// 最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态
		MaxIdle: 10,
		// 最大的激活连接数，表示同时最多有N个连接 ，为0事表示没有限制
		MaxActive: 100,
		//最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
		IdleTimeout: 240 * time.Second,
		// 当链接数达到最大后是否阻塞，如果不的话，达到最大后返回错误
		Wait: true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(url)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return pool, nil
}

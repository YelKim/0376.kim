package redix

import (
	"conf"
	"github.com/garyburd/redigo/redis"
	"github.com/golang/glog"
	"time"
)

type RedisPool struct {
	RedisPool *redis.Pool
}

var redixPool *RedisPool

func GetRedisConn() redis.Conn {
	if redixPool == nil {
		//连接redis
		redixPool = &RedisPool{}
		redixPool.Conn()
	}
	return redixPool.RedisPool.Get()
}

func (this *RedisPool) Conn () {
	this.RedisPool = &redis.Pool{
		MaxActive: 100,
		MaxIdle: 5,
		IdleTimeout: 180 * time.Second,
		Wait: true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.GetConfig().Redis.Conn, redis.DialConnectTimeout(time.Duration(5) * time.Second), redis.DialReadTimeout(time.Duration(10) * time.Second), redis.DialWriteTimeout(time.Duration(10) * time.Second))
			if err != nil {
				glog.Infof("Connect to redis error", err)
				return nil, err
			}
			//验证redis 是否有密码
			if conf.GetConfig().Redis.Passwd != "" {
				if _, err := c.Do("AUTH", conf.GetConfig().Redis.Passwd); err != nil {
					glog.Infof("Connect to redis AUTH error", err)
					c.Close()
					return nil,err
				}
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}
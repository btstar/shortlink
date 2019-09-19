package gredis

import (
	"encoding/json"
	"github.com/fonzie1006/shortlink/pkg/setting"
	"github.com/gomodule/redigo/redis"
	"time"
)

var RedisConn *redis.Pool

func Setup() error {
	RedisConn = &redis.Pool{
		MaxIdle:     setting.RedisSetting.MaxIdle,
		MaxActive:   setting.RedisSetting.MaxActive,
		IdleTimeout: setting.RedisSetting.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", setting.RedisSetting.Host)
			if err != nil {
				return nil, err
			}

			if setting.RedisSetting.Password != "" {
				if _, err := c.Do("AUTH", setting.RedisSetting.Password); err != nil {
					c.Close()
					return nil, err
				}
			}

			if setting.RedisSetting.DB != 0 {
				if _, err := c.Do("SELECT", setting.RedisSetting.DB); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				return err
			}
			return nil
		},
	}
	return nil
}

func redisTest() {
	conn := RedisConn.Get()
	defer conn.Close()

	RedisConn.TestOnBorrow(conn, time.Time{})
}

func Set(key string, data interface{}, t int) error {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	if t > 0 {
		_, err = conn.Do("EXPIRE", key, t)
		if err != nil {
			return err
		}
	}
	return nil

}

func Exists(key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists

}

func Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func GetInt(key string) (int64, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Int64(conn.Do("GET", key))
}

func Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

func LikeDeletes(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err := Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}

func Incr(key string) (int64, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := redis.Int64(conn.Do("INCR", key))
	if err != nil {
		return 0, err
	}

	return value, nil
}

func ResetIncr(key string, data interface{}) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	v, err := redis.Bytes(conn.Do("GETSET", value))
	if err != nil {
		return nil, err
	}

	return v, nil
}

func LikeResetIncr(key string, data interface{}) error {
	conn := RedisConn.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err := ResetIncr(key, data)
		if err != nil {
			return err
		}
	}

	return nil
}

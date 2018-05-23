package router

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

var (
	// 定义常量
	RedisClient *redis.Pool
	REDIS_HOST  string
	REDIS_DB    int
)

func init() {
	// 从配置文件获取redis的ip以及db
	REDIS_HOST = "127.0.0.1:6379"
	REDIS_DB = 1
	// 建立连接池
	RedisClient = &redis.Pool{
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		MaxIdle:     3,
		MaxActive:   10,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", REDIS_HOST)
			if err != nil {
				return nil, err
			}
			/*if _, err := c.Do("AUTH", "2018rds8012"); err != nil {
			    c.Close()
			    return nil, err
			}*/
			// 选择db
			c.Do("SELECT", REDIS_DB)
			return c, nil
		},
	}
}
func RedisGet(key string) (string, error) {
	c := RedisClient.Get()
	defer c.Close()

	v, err := redis.String(c.Do("GET", key))
	if err != nil {
		return "", err
	} else {
		return v, nil
	}
}

func RedisSet(key string, val []byte) error {
	c := RedisClient.Get()
	defer c.Close()
	_, err := c.Do("SET", key, string(val))
	return err
}

func RedisHset(key string, vals map[string]string) error {
	c := RedisClient.Get()
	defer c.Close()
	for k, v := range vals {
		c.Send("HSET", k, v)
	}
	return c.Flush()
	//_,err:=c.DO("HSET")
}

/**
*
 */
func RedisHget(key string, field string) (string, error) {
	c := RedisClient.Get()
	defer c.Close()
	v, err := redis.String(c.Do("HGET", key, field))
	if err != nil {
		return "", err
	} else {
		return v, nil
	}
}

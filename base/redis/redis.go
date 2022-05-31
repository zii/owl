package redis

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Config struct {
	Name         string `json:"name"`
	Addr         string `json:"addr"`
	Active       int    `json:"active"`
	Idle         int    `json:"idle"`
	DialTimeout  int    `json:"dial_timeout"`  // millisecond
	ReadTimeout  int    `json:"read_timeout"`  // millisecond
	WriteTimeout int    `json:"write_timeout"` // millisecond
	IdleTimeout  int    `json:"idle_timeout"`  // second

	DBNum    int    `json:"dbnum"`
	Password string `json:"password"`
}

var redisPoolMap = make(map[string]*redis.Pool)

func Install(configs []*Config) {
	for i := range configs {
		config := configs[i]
		pool := NewRedisPool(config)
		redisPoolMap[config.Name] = pool
	}
}

func MustGetPool(redisName string) *redis.Pool {
	pool, ok := redisPoolMap[redisName]
	if !ok {
		panic("GetRedisClient - Not found pool")
	}
	return pool
}

func GetPool(redisName string) *redis.Pool {
	return redisPoolMap[redisName]
}

func GetRedisPoolClient(redisName string) (conn redis.Conn) {
	pool, ok := redisPoolMap[redisName]
	if !ok {
		panic("GetRedisClient - Not found client")
	}
	conn = pool.Get()
	return
}

func NewRedisPool(c *Config) *redis.Pool {
	cnop := redis.DialConnectTimeout(time.Duration(c.DialTimeout) * time.Millisecond)
	rdop := redis.DialReadTimeout(time.Duration(c.ReadTimeout) * time.Millisecond)
	wrop := redis.DialWriteTimeout(time.Duration(c.WriteTimeout) * time.Millisecond)

	dialFunc := func() (redis.Conn, error) {
		st := time.Now()
		rconn, err := redis.Dial("tcp", c.Addr, cnop, rdop, wrop)
		if err != nil {
			log.Println("redis.Dial took:", time.Since(st), c.Addr)
			return nil, err
		}
		took := time.Since(st)
		if took.Milliseconds() > 1000 {
			log.Println("redis dial slowly:", took)
		}

		if c.Password != "" {
			if _, err = rconn.Do("AUTH", c.Password); err != nil {
				log.Println("Err:", err)
				_ = rconn.Close()
				return nil, err
			}
		}

		_, err = rconn.Do("SELECT", c.DBNum)
		if err != nil {
			_ = rconn.Close()
			return nil, err
		}
		return rconn, nil
	}

	pool := &redis.Pool{
		MaxActive:   c.Active,
		MaxIdle:     c.Idle,
		IdleTimeout: time.Duration(c.IdleTimeout) * time.Second,
		Dial:        dialFunc,
	}
	return pool
}

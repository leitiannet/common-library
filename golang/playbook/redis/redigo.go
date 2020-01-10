/*
单机模式
*/
package main

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

var redisServer = "192.168.32.56:6379"
var redisPassword = "123456"

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func main() {
	pool := newPool(redisServer, redisPassword)
	conn := pool.Get()
	defer conn.Close()

	key := "test_redis"
	value := "this is value for redis"
	_, err := conn.Do("set", key, value)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	reply, err := redis.String(conn.Do("get", key))
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("reply:%v\n", reply)
}

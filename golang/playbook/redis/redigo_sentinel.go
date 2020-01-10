/*
哨兵模式
*/
package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/FZambia/go-sentinel"
	"github.com/garyburd/redigo/redis"
)

var redisServer = "192.168.32.56:26379,192.168.32.57:26379"
var redisPassword = "123456"
var redsiMasterName = "mymaster"

func newSentinelPool(server, password, masterName string) *redis.Pool {
	redisAddrs := strings.Split(server, ",")
	sntnl := &sentinel.Sentinel{
		Addrs:      redisAddrs,
		MasterName: masterName,
		Dial: func(addr string) (redis.Conn, error) {
			timeout := 500 * time.Millisecond
			c, err := redis.DialTimeout("tcp", addr, timeout, timeout, timeout)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
	return &redis.Pool{
		MaxIdle:     3,
		MaxActive:   64,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			masterAddr, err := sntnl.MasterAddr()
			if err != nil {
				return nil, err
			}
			c, err := redis.Dial("tcp", masterAddr)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if !sentinel.TestRole(c, "master") {
				return fmt.Errorf("Role check failed")
			} else {
				return nil
			}
		},
	}
}

func main() {
	pool := newSentinelPool(redisServer, redisPassword, redsiMasterName)
	conn := pool.Get()
	defer conn.Close()

	key := "test_sentinel"
	value := "this is value for sentinel"
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

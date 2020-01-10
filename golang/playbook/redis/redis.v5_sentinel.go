package main

import (
	"fmt"
	"strings"

	"gopkg.in/redis.v5"
)

func main() {
	masterName := "mymaster"
	redisAddr := "192.168.32.56:26379,192.168.32.57:26379"
	redisAddrs := strings.Split(redisAddr, ",")

	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    masterName,
		SentinelAddrs: redisAddrs,
		Password:      "123456",
	})

	key := "test_redis.v5_sentinel"
	value := "this is value for redis.v5 sentinel"
	client.Set(key, value, 0)
	reply, err := client.Get(key).Result()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("reply:%v\n", reply)
}

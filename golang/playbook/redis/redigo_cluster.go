package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/chasex/redis-go-cluster"
)

var redisServer = "192.168.32.45:7000,192.168.32.46:7000"

func newCluster(server string) (*redis.Cluster, error) {
	redisAddrs := strings.Split(server, ",")
	return redis.NewCluster(
		&redis.Options{
			StartNodes:   redisAddrs,
			ConnTimeout:  50 * time.Millisecond,
			ReadTimeout:  50 * time.Millisecond,
			WriteTimeout: 50 * time.Millisecond,
			KeepAlive:    16,
			AliveTime:    60 * time.Second,
		})
}

func main() {
	cluster, err := newCluster(redisServer)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	key := "test-go-cluster"
	value := "this is value for go cluster"
	_, err = cluster.Do("set", key, value)
	if err != nil {
		fmt.Printf("redis.New error: %s\n", err.Error())
		return
	}

	reply, err := redis.String(cluster.Do("get", key))
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("reply:%v\n", reply)
}

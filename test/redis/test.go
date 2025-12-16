package main

import (
	"fmt"
	"time"

	dao_redis "github.com/kongshui/danmu/dao/redis"

	conf "github.com/kongshui/danmu/conf/web"
)

var (
	rdb    = dao_redis.GetRedisClient(config.Redis.Addr, config.Redis.Password, 0, config.Redis.IsCluster, false)
	config = conf.GetConf()
)

func main() {
	t := time.NewTicker(time.Second * 1)
	count := 0
	for {
		<-t.C
		fmt.Println(rdb.Ping())
		count++
		if count == 3 {
			break
		}
	}
}

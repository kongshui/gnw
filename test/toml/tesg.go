package main

import (
	"fmt"

	conf "github.com/kongshui/danmu/conf/node"
)

func main() {
	config := conf.GetConf()
	fmt.Println(config)
}

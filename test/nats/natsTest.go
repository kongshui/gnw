package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	opts := nats.GetDefaultOptions()
	opts.Timeout = 10 * time.Second
	opts.MaxReconnect = 10
	opts.MaxPingsOut = 10
	opts.Servers = []string{"nats://localhost:4222"}

	// 连接到NATS服务器
	// nc, err := nats.Connect("nats://localhost:4222")
	nc, err := opts.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer nc.Close()
	nc.Opts.Dialer.KeepAliveConfig.Enable = true
	fmt.Println("Connected to NATS server")
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()
	// 订阅主题
	go func() {
		d, err := nc.Subscribe("foo", func(msg *nats.Msg) {
			fmt.Printf("Received a message: %s\n", string(msg.Data))
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(3 * time.Second)
		fmt.Println(d.Unsubscribe())
		// <-ctx.Done()
		// log.Println("end")
	}()

	// 发布消息到主题
	err = nc.Publish("foo", []byte("Hello from NATS Go client!"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Published a message to foo")
	go func() {
		t := time.NewTicker(1 * time.Millisecond)
		count := 0
		for {
			<-t.C
			err = nc.Publish("foo", []byte("Hello from NATS Go client!"+strconv.Itoa(count)))
			if err != nil {
				fmt.Println(err)
				return
			}
			count++
		}
	}()
	time.Sleep(10 * time.Second)
	// // 订阅到特定主题
	// _, err = nc.QueueSubscribe("foo", "bar", func(msg *nats.Msg) {
	// 	fmt.Printf("Received a message: %s\n", string(msg.Data))
	// })
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
}

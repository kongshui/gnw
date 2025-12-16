package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

func main() {
	type myKey string

	ctx := context.WithValue(context.Background(), myKey("gender"), "2231")
	fmt.Println(ctx.Value(myKey("gender"))) // 输出: 2231
	os.Open("")
	aaa, _ := os.Open("")
	// aaa = nil
	aaa.Write([]byte("123"))
	for range 10 {
		time.Sleep(10 * time.Second)
	}
}

package main

import (
	"fmt"
	"time"
)

func main() {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
		"f": 6,
		"g": 7,
		"h": 8,
		"i": 9,
	}

	for range 20 {
		count := 0
		for k, v := range m {
			count++
			fmt.Println(v, k)
			if count == 1 {
				break
			}
		}
	}
	Test()
}

func Test() {
	t := time.NewTicker(3 * time.Microsecond)
	defer t.Stop()
	count := 0
	for {
		count++
		<-t.C
		for i := range 10 {
			fmt.Println(i, ":", count)
			time.Sleep(1 * time.Second)
		}
	}
}

package main

import (
	"fmt"
	"sync"
)

var (
	errorPool sync.Pool = sync.Pool{
		New: func() any {
			var err error
			return &err
		},
	}
)

func main() {
	err := errorPool.Get().(*error)
	defer errorPool.Put(err)
	fmt.Println(*err)
}

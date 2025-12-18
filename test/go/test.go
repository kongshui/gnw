package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {
	ctx := context.Background()
	pctx, cancel := context.WithCancel(ctx)
	pctx1, cancel2 := context.WithCancel(ctx)
	defer cancel()
	defer cancel2()
	go func() {
		t := time.NewTicker(2 * time.Second)
		defer t.Stop()
		count := 0
		for range t.C {
			count++
			log.Println("ticker 2s")
			if count >= 5 {
				cancel()
			}
		}
	}()
	select {
	case <-pctx.Done():
		log.Println("pctx canceled")
	case <-pctx1.Done():
		log.Println("pctx1 canceled")
	}
	log.Println("main end")
}

type Test struct {
	Id   int
	Lock sync.Mutex
}

func (t *Test) LockTest(wg *sync.WaitGroup) {
	var err error = fmt.Errorf("lock test error")
	t.Lock.Lock()

	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("lock test")
		wg.Done()
	}()
	t.Lock.Unlock()
	if err != nil {
		log.Println("lock test error:" + err.Error())
		return
	}
}

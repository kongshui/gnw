package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {
	var (
		t  Test
		wg sync.WaitGroup
	)
	for range 100 {
		wg.Add(1)
		t.LockTest(&wg)
	}
	wg.Wait()
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

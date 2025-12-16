package main

import (
	"fmt"
	"log"
)

func TestRecover() {
	var isok bool = false
	defer TestRecover2(isok)
	isok = true
	panic("TestRecover")
}

func TestRecover2(ok bool) {
	if err := recover(); err != nil {
		if fmt.Sprintf("%v", err) == "TestRecover" {
			log.Println("111111")
		}
		log.Println("TestRecover2 recover err:", err, "ok:", ok)
	}
}

func main() {
	TestRecover()
}

package main

import "fmt"

type Node struct {
	Val []int
}

func (n *Node) Add(i int) {
	n.Val = append(n.Val, i)
}

func (n *Node) Read() {
	for _, v := range n.Val {
		println(v)
	}
}

func main() {
	aaa := make(map[string][]string)
	aaa["a"] = []string{"a", "b", "c"}
	fmt.Println(aaa)
}

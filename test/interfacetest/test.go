package main

import (
	"fmt"
	"net"
	"reflect"

	"github.com/gorilla/websocket"
)

func main() {
	var d D
	var c *net.Conn
	var e *websocket.Conn
	fmt.Println(reflect.TypeOf(c))
	fmt.Println(reflect.TypeOf(&d))
	if reflect.TypeOf(e).String() == "*websocket.Conn" {
		fmt.Println("true")
	}
	if reflect.TypeOf(c).String() == "*net.Conn" {
		fmt.Println("true")
	}
}
func Test_GetConn(d any) {
	fmt.Println(d.(*net.Conn))
	fmt.Println(reflect.TypeOf(d))
}

type (
	Test interface {
		SetNodeType(t any)
		GetNodeType() any
	}
	D struct {
		node_type int
	}
)

func Test_GetNodeType(d Test) {
}

// 设置节点类型
func (d *D) SetNodeType(v any) {
	d.node_type = v.(int)
}

func (d *D) GetNodeType() any {
	return d.node_type
}
func (d *D) CCCCC() any {
	return d.node_type
}

func Test_SetNodeType(d Test, t int) {
	d.SetNodeType(t)
	fmt.Println(d.GetNodeType())
}

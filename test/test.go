package main

import (
	"fmt"
	"sync"
)

type Test struct {
	Name string
	Age  int
	Lock *sync.RWMutex
}

func (t *Test) SetName(name string) bool {
	t.Lock.RLock()
	defer t.Lock.RUnlock()
	t.Lock.Lock()
	defer t.Lock.Unlock()
	t.Name = name
	return t.Name == name
}

func main() {
	// var t *Test = &Test{Name: "test", Age: 1, Lock: &sync.RWMutex{}}
	// t.SetName("test2")
	// fmt.Println(t.Name)
	aaa := "中国人"
	fmt.Println(len(aaa))
	fmt.Println(len([]byte(aaa)))
	//______________________________________________________________
	// t := &Test{Name: "test", Age: 1}
	// fmt.Println(GetShape(t))
	// fmt.Println(t.aaa)
	// var wg msg.MsgConn
	// wg.SetGroupId("test3")
	// fmt.Println(wg.GetGroupId())
	// t := &Test{Name: "test", Age: 1}
	// msgIdB := make([]byte, 4)
	// fmt.Println(len(msgIdB))
	// binary.BigEndian.PutUint32(msgIdB, 888888888)
	// d, _ := json.Marshal(t)
	// fmt.Println(len(d))
	// aaa := make([]byte, 0)
	// aaa = append(aaa, msgIdB...)
	// aaa = append(aaa, d...)
	// fmt.Println(len(aaa))
	// c := Test{}
	// if err := json.Unmarshal(aaa[4:], &c); err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(c)
	//______________________________________________________________
	// a := map[uint32]func([]byte, msginterface.MsgConn, []byte){
	// 	1: func(id []byte, conn msginterface.MsgConn, data []byte) {
	// 		c := Test{}
	// 		if err := json.Unmarshal(data, &c); err != nil {
	// 			fmt.Println(err)
	// 		}
	// 		fmt.Println(c, 11111111)
	// 	},
	// }
	// ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	// defer cancel()
	// l, _ := net.Listen("tcp", ":8888")
	// c, _ := l.Accept()
	// defer c.Close()
	// conn := msg.NewMsgConn(c)
	// conn.SetOnline(true)
	// conn.SetCtx(ctx)
	// conn.SetCancel(cancel)
	// go conn.ReceiveMessage(a)
	// go conn.Keepalived()
	// for i := range 15 {
	// 	_, err := conn.MessageWrite([]byte(""))
	// 	if err != nil {
	// 		fmt.Println(err, 11111)
	// 		break
	// 	}
	// 	time.Sleep(1 * time.Second)
	// 	fmt.Println(i)
	// }
	// fmt.Println("end")
	// go func() {
	// 	for {
	// 		time.Sleep(1 * time.Second)
	// 		fmt.Println(conn.GetConn())
	// 	}
	// }()
	// <-ctx.Done()
}

func (t *Test) GetName() string {
	return t.Name
}
func (t Test) Close() error {
	return nil
}

type (
	MessageConn interface {
		// GetGroup() string
		GetName() string
		SetName(string) bool
		// net.Conn
	}
)

// 接口使用
func SetShape(s MessageConn, name string) {
	s.SetName(name)
}

// 接口使用
func GetShape(s MessageConn) string {
	return s.GetName()
}

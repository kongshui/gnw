package main

import (
	"fmt"

	"github.com/google/uuid"
)

// type Test struct {
// 	Conn net.Conn
// 	uid  uuid.UUID
// }

// var (
// 	testUid map[uuid.UUID]Test = make(map[uuid.UUID]Test)
// )

func main() {
	uid := uuid.New()
	uidByte, _ := uid.MarshalBinary()
	uidString := uid.String()
	uidByteString := string(uidByte)
	fmt.Println(uidByteString, uidString, uidByte)
	a := []byte(uidByteString)
	b, _ := uuid.FromBytes(a)
	fmt.Println(b.String(), a)
	// fmt.Println(os.Environ())
	var s []string = make([]string, 0)
	Test(&s)
	fmt.Println(s)
	// l, _ := net.Listen("tcp", "127.0.0.1:8080")
	// count := 0
	// for {
	// 	count++
	// 	e, _ := l.Accept()
	// 	uid = uuid.New()
	// 	testUid[uid] = Test{Conn: e, uid: uid}
	// 	if e == testUid[uid].Conn {
	// 		fmt.Println("true")
	// 	}
	// 	if count == 2 {
	// 		break
	// 	}
	// }

}

func Test(a *[]string) {
	*a = append(*a, "test")
}

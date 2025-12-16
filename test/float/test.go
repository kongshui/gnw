package main

import (
	"fmt"
	"log"
	"path"
	"strconv"
	"strings"
)

func main() {
	var a float64 = 1
	fmt.Println(strconv.FormatFloat(a, 'f', 2, 64))
	fmt.Println(strconv.ParseFloat("3.01", 64))
	b := "/a/b/c"
	fmt.Println(strings.Split(b, "/")[0], 1, strings.Split(b, "/")[1], 2, strings.Split(b, "/")[2], 3, strings.Split(b, "/")[3])

	log.Println(path.Join("//", "bbb", "ccc"))
}

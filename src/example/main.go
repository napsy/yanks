package main

import (
	"fmt"
	"time"
	_ "yanks"
)

var i int
var m map[int][]byte = make(map[int][]byte)

func allocs() {
	f := make([]byte, 2000)
	f[222] = 'f'
	m[i] = f
	i++
}
func main() {
	fmt.Println("vim-go")
	for {
		time.Sleep(time.Second)
		allocs()
	}
	time.Sleep(time.Minute)
}

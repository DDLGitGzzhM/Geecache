package main

import (
	"fmt"
	"sync"
)

var m sync.Mutex
var set = make(map[int]bool, 0)
var wg sync.WaitGroup

func printOnce(num int) {
	m.Lock()
	if _, ok := set[num]; !ok {
		fmt.Println(num)
	}
	set[num] = true
	m.Unlock()
	defer wg.Done()
}

func main() {

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go printOnce(100)
	}
	wg.Wait()
}

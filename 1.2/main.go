package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var buffer = -1
var cond = sync.NewCond(&sync.Mutex{})
var cond2 = sync.NewCond(&sync.Mutex{})
var wg sync.WaitGroup

func consume(id int) int {
	cond2.Signal()

	cond.L.Lock()
	wg.Add(1)
	for buffer == -1 {
		cond.Wait()
	}
	value := buffer
	wg.Done()
	cond.L.Unlock()
	return value
}

func produce(n int) {
	cond.L.Lock()
	buffer = n
	fmt.Printf("--- broadcasted %d\n", n)
	cond.L.Unlock()
	cond.Broadcast()
	wg.Wait()
	cond.L.Lock()
	buffer = -1
	wg = sync.WaitGroup{}
	cond.L.Unlock()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	consumers := 2
	for i := 0; i < consumers; i++ {
		go func(id int) {
			for {
				fmt.Printf("Consumed: %d\n", consume(id))
			}
		}(i)
	}

	for {
		produce(rand.Intn(100))

		cond2.L.Lock()
		cond2.Wait()
		cond2.L.Unlock()
	}
}

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type consumerRes struct {
	id   int
	x    *int
	cond *sync.Cond
}

var exit = make(chan bool)

var consQueueCond = sync.NewCond(&sync.Mutex{})
var consQueue = make([]*consumerRes, 0)

func freeConsumers(n int) {
	consQueueCond.L.Lock()
	for len(consQueue) == 0 {
		consQueueCond.Wait()
	}
	for _, cons := range consQueue {
		cons.cond.L.Lock()
		cons.x = &n
		fmt.Printf("producer signaled for consumer %d\n", cons.id)
		cons.cond.Signal()
		cons.cond.L.Unlock()
	}
	consQueue = make([]*consumerRes, 0)
	consQueueCond.L.Unlock()
}

func consumerThread(id int) {
	for {
		cond := sync.NewCond(&sync.Mutex{})
		res := consumerRes{id, nil, cond}

		consQueueCond.L.Lock()
		consQueue = append(consQueue, &res)
		consQueueCond.Signal()
		consQueueCond.L.Unlock()

		cond.L.Lock()
		for res.x == nil {
			cond.Wait()
		}
		cond.L.Unlock()

		fmt.Printf("%d consumer, consumed: %d\n", id, *res.x)
	}
}

func producerThread() {
	// time.Sleep(2 * time.Second)
	// for i := 0; i < 10; i++ {
	for {
		x := rand.Intn(100)
		fmt.Printf("produced %d\n", x)

		freeConsumers(x)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	consumers := 1 + rand.Intn(10)
	for i := 0; i < consumers; i++ {
		go consumerThread(i + 1)
	}

	go producerThread()

	<-exit
}

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type consumerDocument struct {
	x    *int
	cond *sync.Cond
}

var exit = make(chan bool)

var queue = NewThreadSafeQueue(10)

func consume() int {
	doc := newConsumerDocument()
	queue.Push(doc)

	doc.waitFilling()

	return *doc.x
}

func produce(n int) {
	queue.Lock()
	for queue.Empty() {
		queue.Wait()
	}
	for !queue.Empty() {
		doc := toDocument(queue.UnsafePop())
		doc.fill(n)
	}
	queue.Unlock()
}

func consumerThread(id int) {
	for {
		value := consume()

		fmt.Printf("%d consumer, consumed: %d\n", id, value)
	}
}

func producerThread() {
	// time.Sleep(2 * time.Second)
	// for i := 0; i < 10; i++ {
	for {
		x := rand.Intn(100)
		fmt.Printf("produced %d\n", x)

		produce(x)
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

func toDocument(value interface{}, err error) *consumerDocument {
	return value.(*consumerDocument)
}

func newConsumerDocument() *consumerDocument {
	return &consumerDocument{
		nil,
		sync.NewCond(&sync.Mutex{}),
	}
}

func (doc *consumerDocument) waitFilling() {
	doc.cond.L.Lock()
	for doc.x == nil {
		doc.cond.Wait()
	}
	doc.cond.L.Unlock()
}

func (doc *consumerDocument) fill(n int) {
	doc.cond.L.Lock()
	doc.x = &n
	doc.cond.Signal()
	doc.cond.L.Unlock()
}

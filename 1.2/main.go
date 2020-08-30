package main

import (
	"fmt"
	"math/rand"
	"time"
)

var buffer = make(chan (chan int), 10)

func consume() int {
	doc := make(chan int)
	buffer <- doc
	return <-doc
}

func produce(n int) {
	if consumers := len(buffer); consumers > 0 {
		for ; consumers > 0; consumers-- {
			doc := <-buffer
			doc <- n
		}
	} else {
		doc := <-buffer
		doc <- n
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	consumers := 1 + rand.Intn(10)
	for i := 0; i < consumers; i++ {
		go func() {
			for {
				fmt.Printf("consumed: %d\n", consume())
			}
		}()
	}

	for {
		n := rand.Intn(100)
		fmt.Printf("produced: %d\n", n)
		produce(n)
	}
}

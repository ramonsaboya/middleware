package main

import (
	"fmt"
	"math/rand"
	"time"
)

type bridge chan bool

type side struct {
	displayName string
	cars        chan string
}

var exit = make(chan bool)

var carColors = [...]string{"blue", "green", "red", "yellow"}

func randomColor() string {
	return carColors[rand.Intn(len(carColors))]
}

func randomDuration(min, max int) time.Duration {
	seconds := time.Duration(min + rand.Intn(max-min+1))
	return seconds * time.Second
}

func carProducer(side side) {
	for {
		time.Sleep(randomDuration(3, 5))
		for i := 0; i < rand.Intn(3)+1; i++ {
			color := randomColor()
			fmt.Printf("%s car spawned on %s side\n", color, side.displayName)
			side.cars <- color
		}
		fmt.Println()
	}
}

func crossCars(commonBridge bridge, side side) {
	for {
		<-commonBridge
		crossing := true
		for crossing {
			select {
			case color, _ := <-side.cars:
				fmt.Printf("----- %s car crossed from the %s\n", color, side.displayName)
			default:
				crossing = false
				fmt.Println()
			}
		}
		commonBridge <- true
		time.Sleep(randomDuration(1, 2))
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	leftSide := side{"left", make(chan string, 10)}
	rightSide := side{"right", make(chan string, 10)}
	go carProducer(leftSide)
	go carProducer(rightSide)

	commonBridge := make(bridge)
	go crossCars(commonBridge, leftSide)
	go crossCars(commonBridge, rightSide)

	time.Sleep(randomDuration(5, 6))
	commonBridge <- true

	<-exit
}

// Seán Rourke
// C00251168

package main

import (
	"fmt"
	"time"
)

func producer(theChan chan int) {
	for i := 0; i < 5; i++ { // loop for producers
		time.Sleep(time.Second)
		fmt.Println("Producer: producing", i)
		theChan <- i // add to channel
	}
	close(theChan) // close channel when all values are added
}

func consumer(theChan <-chan int) {
	for i := range theChan { // loop for consumers
		time.Sleep(time.Second)
		fmt.Println("Consumer: consuming", i)
	}
}

func main() {
	theChan := make(chan int)
	go producer(theChan)
	go consumer(theChan)
	time.Sleep(time.Second * 10) // sleep to allow processes to run
}

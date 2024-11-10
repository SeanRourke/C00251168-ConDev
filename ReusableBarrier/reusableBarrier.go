// Seán Rourke
// C00251168

package main

import (
	"fmt"
	"sync"
	"time"
)

func doStuff(goNum int, arrived *int, max int, wg *sync.WaitGroup, sharedLock *sync.Mutex, theChan chan bool) bool {
	for j := 0; j < 5; j++ { // loop for number of times barrier is used
		time.Sleep(time.Second)
		fmt.Println("Part A", goNum)
		sharedLock.Lock()
		*arrived++           // count for all Part A processes
		if *arrived == max { //last to arrive, signal others to go
			sharedLock.Unlock()
			theChan <- true
			<-theChan
		} else { // not all here yet we wait until signal
			sharedLock.Unlock()
			<-theChan
			theChan <- true
		}
		sharedLock.Lock()
		*arrived--
		if *arrived == 0 {
			sharedLock.Unlock()
			theChan <- true
			<-theChan
		} else {
			sharedLock.Unlock()
			<-theChan
			theChan <- true
		}
		fmt.Println("PartB", goNum) // all Part A arrived, proceed with Part B
	}
	wg.Done()
	return true
}

func main() {
	totalRoutines := 10
	arrived := 0
	var wg sync.WaitGroup
	wg.Add(totalRoutines)
	var theLock sync.Mutex
	theChan := make(chan bool)
	for i := range totalRoutines {
		go doStuff(i, &arrived, totalRoutines, &wg, &theLock, theChan)
	}
	wg.Wait() //wait for everything to finish before exiting
}

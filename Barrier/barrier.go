// Seán Rourke
// C00251168

package main

import (
	"fmt"
	"sync"
	"time"
)

// Place a barrier in this function --use Mutex's and Semaphores
func doStuff(goNum int, arrived *int, max int, wg *sync.WaitGroup, sharedLock *sync.Mutex, theChan chan bool) bool {
	time.Sleep(time.Second)
	fmt.Println("Part A", goNum)
	//we wait here until everyone has completed part A
	sharedLock.Lock()
	*arrived++
	if *arrived == max { //last to arrive -signal others to go
		sharedLock.Unlock() //unlock before any potentially blocking code
		theChan <- true
		<-theChan
	} else { //not all here yet we wait until signal
		sharedLock.Unlock() //unlock before any potentially blocking code
		<-theChan
		theChan <- true //once we get through send signal to next routine to continue
	} //end of if-else
	//we wait here until everyone has completed part A
	fmt.Println("PartB", goNum)
	wg.Done()
	return true
}

func main() {
	totalRoutines := 10
	arrived := 0
	var wg sync.WaitGroup
	wg.Add(totalRoutines)
	//we will need some of these
	theChan := make(chan bool) //use unbuffered channel in place of semaphore
	var theLock sync.Mutex
	for i := range totalRoutines { //create the go Routines here
		go doStuff(i, &arrived, totalRoutines, &wg, &theLock, theChan)
	}

	wg.Wait() //wait for everyone to finish before exiting
}

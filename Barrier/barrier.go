// Seán Rourke
// C00251168

package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

// Place a barrier in this function --use Mutex's and Semaphores
func doStuff(goNum int, arrived *int, max int, wg *sync.WaitGroup, sharedLock *sync.Mutex, theSem *semaphore.Weighted, ctx context.Context) bool {
	time.Sleep(time.Second)
	fmt.Println("Part A", goNum)
	sharedLock.Lock()
	*arrived++
	if *arrived == max {
		theSem.Release(1)
		sharedLock.Unlock()
		theSem.Acquire(ctx, 1)
	} else {
		sharedLock.Unlock()
		theSem.Acquire(ctx, 1)
		theSem.Release(1)
	}
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
	ctx := context.TODO()
	var theLock sync.Mutex
	sem := semaphore.NewWeighted(0)
	for i := range totalRoutines { //create the go Routines here
		go doStuff(i, &arrived, totalRoutines, &wg, &theLock, sem, ctx)
	}

	wg.Wait() //wait for everyone to finish before exiting
}

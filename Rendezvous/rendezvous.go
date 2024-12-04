// Seán Rourke
// C00251168

package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

func rendezvous(wg *sync.WaitGroup, num int) bool {
	var x time.Duration
	x = time.Duration(rand.IntN(5))
	time.Sleep(x * time.Second)
	fmt.Println("Part A", num)

	fmt.Println("Part B", num)
	wg.Done()
	return true
}

func main() {
	var wg sync.WaitGroup
	threads := 5

	wg.Add(threads)
	for thread := range threads {
		go rendezvous(&wg, thread)
	}
	wg.Wait()
}

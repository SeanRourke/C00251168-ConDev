package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Philosopher struct {
	Id        int
	LeftFork  *sync.Mutex
	RightFork *sync.Mutex
}

const (
	RandSecond      = 1e9
	NOfPhilosophers = 5
	Phil            = "Philosopher"
)

func main() {
	counter := make(chan int, 1)
	counter <- 0

	// Create a fork (mutex) for each philosopher
	forks := make([]*sync.Mutex, NOfPhilosophers)
	for i := 0; i < NOfPhilosophers; i++ {
		forks[i] = &sync.Mutex{}
	}

	// Create a philosopher for each position at the table
	philosophers := make([]*Philosopher, NOfPhilosophers)
	for i := 0; i < NOfPhilosophers; i++ {
		philosophers[i] = &Philosopher{
			Id:        i + 1,
			LeftFork:  forks[i],
			RightFork: forks[(i+1)%NOfPhilosophers],
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(NOfPhilosophers)

	fmt.Printf("There are %v philosophers sitting at a table\n", NOfPhilosophers)
	for _, phil := range philosophers {
		go func(syncer *sync.WaitGroup, ph *Philosopher) {
			defer syncer.Done()
			ph.dining(counter)
			fmt.Printf("%s %v - is done dining\n", Phil, ph.Id)
		}(&wg, phil)
	}
	wg.Wait()
	c := <-counter
	fmt.Printf("%v philosophers finished eating!\n", c)
}

// dining process: philosopher tries to acquire forks, eats, then releases forks and increments the counter
func (phil *Philosopher) dining(counter chan int) {
	phil.getForks()
	phil.eating()
	phil.returnForks()
	c := <-counter
	c += 1
	counter <- c
}

// getForks the process of getting the forks
func (phil *Philosopher) getForks() {
	phil.thinking()
	fmt.Printf("%s %v - is trying to get forks\n", Phil, phil.Id)
	// Lock left fork first, then right fork
	phil.LeftFork.Lock()
	fmt.Printf("%s %v - got the left fork\n", Phil, phil.Id)
	phil.RightFork.Lock()
	fmt.Printf("%s %v - got the right fork\n", Phil, phil.Id)
}

// returnForks releases the forks after eating
func (phil *Philosopher) returnForks() {
	// Unlock the forks to release them for others
	phil.LeftFork.Unlock()
	phil.RightFork.Unlock()
	fmt.Printf("%s %v - returned forks\n", Phil, phil.Id)
}

// thinking simulates thinking with a random delay
func (phil *Philosopher) thinking() {
	t := time.Duration(rand.Int63n(RandSecond))
	fmt.Printf("%s %v - is thinking for %v\n", Phil, phil.Id, t)
	time.Sleep(t)
}

// eating simulates eating with a random delay
func (phil *Philosopher) eating() {
	t := time.Duration(rand.Int63n(RandSecond))
	fmt.Printf("%s %v - is eating for %v\n", Phil, phil.Id, t)
	time.Sleep(t)
}

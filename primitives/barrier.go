package primitives

import (
	"fmt"
	"sync"
)

func doWork(mutex, barrier *CondSemaphore, thread int, wg *sync.WaitGroup, arrived *int, workers int) {
	defer wg.Done()
	fmt.Printf("Thread %d starting \n", thread)
	mutex.Wait()
	*arrived += 1
	if *arrived == workers {
		fmt.Printf("unlocking barrier at thread %d \n", thread)
		barrier.Signal()
	}
	mutex.Signal()
	barrier.Wait()
	fmt.Printf("Thread %d working \n", thread)
	barrier.Signal()
}

func main() {
	workers := 15
	arrived := 0
	wg := sync.WaitGroup{}
	for k := 0; k < 3; k++ {
		fmt.Printf("starting iteration %d\n", k)
		wg.Add(workers)
		m := NewCondSemaphore(1)
		barrier := NewCondSemaphore(0)
		for i := 0; i < workers; i++ {
			go doWork(m, barrier, i, &wg, &arrived, workers)
		}
		wg.Wait()
	}
}

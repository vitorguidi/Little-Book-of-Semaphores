package primitives

import (
	"fmt"
	"sync"
)

func prepare() {}

func doWorkReusable(mutex, barrierEntry, barrierLeave *CondSemaphore, thread int, wg *sync.WaitGroup, arrived *int, workers int) {
	defer wg.Done()
	prepare()
	//must unlock enter barrier before coming in here
	//can only enter when all threads have prepared
	mutex.Wait()
	*arrived -= 1
	if *arrived == 0 {
		barrierEntry.Signal()
		barrierLeave.Wait() //we end up with 1 extra signal after all threads did their work and come back here
	}
	mutex.Signal()

	barrierEntry.Wait()
	barrierEntry.Signal()
	fmt.Printf("Thread %d starting \n", thread)
	mutex.Wait()
	*arrived += 1
	if *arrived == workers {
		fmt.Printf("unlocking barrier at thread %d \n", thread)
		barrierLeave.Signal()
		barrierEntry.Wait() //we end up with 1 extra signal when all threads leave prepare
	}
	mutex.Signal()
	barrierLeave.Wait()
	fmt.Printf("Thread %d working \n", thread)
	barrierLeave.Signal()
}

//func main() {
//	workers := 4
//	arrived := workers
//	wg := sync.WaitGroup{}
//	m := NewCondSemaphore(1)
//	barrierEntry := NewCondSemaphore(0)
//	barrierLeave := NewCondSemaphore(1)
//	for k := 0; k < 3; k++ {
//		fmt.Printf("starting iteration %d\n", k)
//		wg.Add(workers)
//		for i := 0; i < workers; i++ {
//			go doWorkReusable(m, barrierEntry, barrierLeave, i, &wg, &arrived, workers)
//		}
//		wg.Wait()
//	}
//}

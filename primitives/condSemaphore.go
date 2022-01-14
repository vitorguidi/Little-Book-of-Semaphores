package primitives

import (
	"sync"
)

type CondSemaphore struct {
	cnt  int
	cond *sync.Cond
}

func NewCondSemaphore(limit int) *CondSemaphore {
	return &CondSemaphore{
		cnt:  limit,
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func (sem *CondSemaphore) Signal() {
	sem.cond.L.Lock()
	sem.cnt += 1
	sem.cond.L.Unlock()
	sem.cond.Signal()
}

func (sem *CondSemaphore) Wait() {
	sem.cond.L.Lock()
	for sem.cnt == 0 {
		sem.cond.Wait()
	}
	sem.cnt--
	sem.cond.L.Unlock()
}

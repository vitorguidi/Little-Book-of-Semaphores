package main

import (
	"fmt"
	"lbs/primitives"
	"sync"
)

var test chan byte

type RWMutex struct {
	people                      int
	roomEmpty, readTurnstile, m *primitives.CondSemaphore
}

func NewRWMutex() *RWMutex {
	return &RWMutex{
		roomEmpty:     primitives.NewCondSemaphore(1),
		readTurnstile: primitives.NewCondSemaphore(1),
		m:             primitives.NewCondSemaphore(1),
	}
}

func (rwm *RWMutex) Wait() {
	rwm.readTurnstile.Wait() //avoids new readers from entering the room
	rwm.roomEmpty.Wait()     //waits for last reader to leave
}

func (rwm *RWMutex) Signal() {
	rwm.readTurnstile.Signal() //allows readers to pass through entry or next writer to snatch
	rwm.roomEmpty.Signal()     //writer leaves room
}

func (rwm *RWMutex) RWait() {
	rwm.readTurnstile.Wait()   //waits for turnstile to open
	rwm.readTurnstile.Signal() //allows next reader to go through or writer to snatch the entry lock
	rwm.incrementLightSwitch() //incs counter of readers inside room. if writer is there, gets blocked by roomEmpty
}

func (rwm *RWMutex) RSignal() {
	rwm.decrementLightSwitch() //reader leaves room
}

func (rwm *RWMutex) incrementLightSwitch() {
	rwm.m.Wait()
	rwm.people += 1
	if rwm.people == 1 {
		rwm.roomEmpty.Wait()
	}
	rwm.m.Signal()
}

func (rwm *RWMutex) decrementLightSwitch() {
	rwm.m.Wait()
	rwm.people -= 1
	if rwm.people == 0 {
		rwm.roomEmpty.Signal()
	}
	rwm.m.Signal()
}

func reader(wg *sync.WaitGroup, i int, m *RWMutex) {
	defer wg.Done()
	m.RWait()
	fmt.Printf("{")
	fmt.Printf("}")
	m.RSignal()
}

func writer(wg *sync.WaitGroup, i int, m *RWMutex) {
	defer wg.Done()
	m.Wait()
	fmt.Printf("(")
	fmt.Printf(")")
	m.Signal()
}

func validator() {
	openKey := 0
	openParens := 0
	ans := true
	for range test {
		letter := <-test
		switch letter {
		case '(':
			openParens += 1
			if openKey > 0 {
				ans = false
			}

		case ')':
			openParens -= 1

		case '{':
			openKey += 1
			if openParens > 0 {
				ans = false
			}

		case '}':
			openKey -= 1

		}
		fmt.Println(ans)
	}
}

func main() {
	test = make(chan byte)
	readers := 1000000
	writers := 1000000
	rwm := NewRWMutex()
	wg := sync.WaitGroup{}
	go validator()
	for i := 0; i < writers; i++ {
		wg.Add(1)
		go writer(&wg, i, rwm)
	}
	for i := 0; i < readers; i++ {
		wg.Add(1)
		go reader(&wg, i, rwm)
	}

	wg.Wait()
	close(test)

}

package primitives

type RWMutex struct {
	people                      int
	roomEmpty, readTurnstile, m *CondSemaphore
}

func NewRWMutex() *RWMutex {
	return &RWMutex{
		roomEmpty:     NewCondSemaphore(1),
		readTurnstile: NewCondSemaphore(1),
		m:             NewCondSemaphore(1),
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

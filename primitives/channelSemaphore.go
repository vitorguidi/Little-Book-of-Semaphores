package primitives

type ChannelSemaphore struct {
	ch chan struct{}
}

func NewChannelSemaphore(limit int) *ChannelSemaphore {
	return &ChannelSemaphore{
		ch: make(chan struct{}, limit),
	}
}

func (sem *ChannelSemaphore) Signal() {
	sem.ch <- struct{}{}
}

func (sem *ChannelSemaphore) Wait() {
	<-sem.ch
}

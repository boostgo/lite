package async

// Semaphore is tool for managing goroutines count at a time.
// Create semaphore - async.NewSemaphore(N). N - num of max goroutines at a time
type Semaphore struct {
	c chan struct{}
}

func (s *Semaphore) Acquire() {
	s.c <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.c
}

func (s *Semaphore) Close() {
	close(s.c)
}

func NewSemaphore(size int) *Semaphore {
	return &Semaphore{
		c: make(chan struct{}, size),
	}
}

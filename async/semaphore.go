package async

// Semaphore is tool for managing goroutines count at a time
type Semaphore struct {
	c chan struct{}
}

// NewSemaphore create Semaphore. size - num of max goroutines at a time
func NewSemaphore(size int) *Semaphore {
	return &Semaphore{
		c: make(chan struct{}, size),
	}
}

// Acquire add one more semaphore to pool.
//
// Here is becoming "wait/hold" moment.
func (s *Semaphore) Acquire() {
	s.c <- struct{}{}
}

// Release remove one semaphore from pool.
//
// It allows pool to acquire one more semaphore
func (s *Semaphore) Release() {
	<-s.c
}

func (s *Semaphore) Close() {
	close(s.c)
}

package health

import (
	"context"
	"github.com/boostgo/lite/errs"
	"golang.org/x/sync/errgroup"
	"time"
)

type session struct {
	parentCtx context.Context
	checkers  []Checker
	timeout   time.Duration
	logging   bool

	statuses []Status
}

func newSession(parentCtx context.Context, checkers []Checker) *session {
	return &session{
		parentCtx: parentCtx,
		checkers:  checkers,

		statuses: make([]Status, 0),
	}
}

func (s *session) Timeout(timeout time.Duration) *session {
	s.timeout = timeout
	return s
}

func (s *session) Logging(logging bool) *session {
	s.logging = logging
	return s
}

func (s *session) Check() []Status {
	wg := errgroup.Group{}
	statusesChan := make(chan Status, len(s.checkers))
	for _, checker := range s.checkers {
		wg.Go(func() error {
			var ctx context.Context
			var cancel context.CancelFunc
			if s.timeout > 0 {
				ctx, cancel = context.WithTimeout(s.parentCtx, s.timeout)
				defer cancel()
			}

			select {
			case <-time.After(s.timeout):
				return errs.ErrTimeout
			case status := <-checker.StatusAsync(ctx):
				statusesChan <- status
				return nil
			}
		})
	}

	_ = wg.Wait()
	close(statusesChan)

	statuses := make([]Status, 0, len(statusesChan))
	for status := range statusesChan {
		statuses = append(statuses, status)
	}

	return statuses
}

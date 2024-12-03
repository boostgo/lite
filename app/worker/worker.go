package worker

import (
	"context"
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/log"
	"github.com/boostgo/lite/system/life"
	"github.com/boostgo/lite/system/trace"
	"github.com/boostgo/lite/system/try"
	"time"
)

// Worker is job/cron based structure.
type Worker struct {
	name         string
	fromStart    bool
	duration     time.Duration
	action       func(ctx context.Context) error
	errorHandler func(error) bool
	stopper      chan bool
	done         chan struct{}
	traceMaster  bool
	timeout      time.Duration
}

// New creates [Worker] object
func New(name string, duration time.Duration, action func(ctx context.Context) error) *Worker {
	return &Worker{
		name:        name,
		duration:    duration,
		action:      action,
		stopper:     make(chan bool),
		done:        make(chan struct{}, 1),
		traceMaster: trace.AmIMaster(),
	}
}

// FromStart sets flag for starting worker from start.
func (worker *Worker) FromStart() *Worker {
	worker.fromStart = true
	return worker
}

// Timeout sets timeout duration for working action timeout.
func (worker *Worker) Timeout(timeout time.Duration) *Worker {
	worker.timeout = timeout
	return worker
}

// ErrorHandler sets custom error handler from action
func (worker *Worker) ErrorHandler(handler func(error) bool) *Worker {
	if handler == nil {
		return worker
	}

	worker.errorHandler = handler
	return worker
}

// runAction runs provided action with context and try function and trace id.
func (worker *Worker) runAction() error {
	ctx := context.Background()
	var cancel context.CancelFunc

	if worker.traceMaster {
		ctx = trace.Set(ctx, trace.String())
	}

	if worker.duration > 0 {
		ctx, cancel = context.WithTimeout(ctx, worker.duration)
		defer cancel()
	}

	return try.Ctx(ctx, worker.action)
}

// Run runs worker with provided duration
func (worker *Worker) Run() {
	logger := log.Namespace("worker")

	if worker.fromStart {
		if err := worker.runAction(); err != nil {
			logger.
				Error().
				Str("worker", worker.name).
				Err(err).
				Msg("Start worker action")
		}
	}

	go func() {
		ticker := time.NewTicker(worker.duration)
		defer ticker.Stop()

		life.Tear(func() error {
			// teardown will make main goroutine wait till worker will not be done
			<-worker.done
			return nil
		})

		for {
			select {
			case <-life.Context().Done():
				logger.Info().Str("worker", worker.name).Msg("Stop worker by context")
				worker.done <- struct{}{}
				return
			case <-worker.stopper:
				logger.Info().Str("worker", worker.name).Msg("Stop worker by stopper")
				worker.done <- struct{}{}
				return
			case <-ticker.C:
				if err := worker.runAction(); err != nil {
					logger.Error().Str("worker", worker.name).Err(err).Msg("Ticker worker action")

					if errs.IsType(err, "Panic") {
						worker.stopper <- true
						continue
					}

					if worker.errorHandler != nil {
						if !worker.errorHandler(err) {
							worker.stopper <- true
							continue
						}
					}
				}
			}
		}
	}()
}

// Run created worker object and runs by itself. It is like "short" version of using [Worker]
func Run(name string, duration time.Duration, action func(ctx context.Context) error, fromStart ...bool) {
	worker := New(name, duration, action)
	if len(fromStart) > 0 && fromStart[0] {
		worker.FromStart()
	}
	worker.Run()
}

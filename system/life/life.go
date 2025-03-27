package life

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"slices"
	"sync"
	"time"
)

var (
	_once sync.Once
	_life *life
)

// life controls the application lifetime.
//
// Contain global app context, context's cancel function and teardown functions
type life struct {
	ctx         context.Context
	cancel      context.CancelFunc
	tears       []func() error
	gracefulLog func()
}

// Cancel global app context
func (life *life) Cancel() {
	life.cancel()
}

func instance() *life {
	_once.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		_life = &life{
			ctx:    ctx,
			cancel: cancel,
			tears:  make([]func() error, 0),
		}
	})

	return _life
}

// Init initialize global life instance
func Init() {
	instance()
}

// Context returns global context
func Context() context.Context {
	return instance().ctx
}

// Cancel call global context cancel function
func Cancel() {
	instance().cancel()
}

// GracefulLog set function which calls when Cancel called
func GracefulLog(gracefulLog func()) {
	instance().gracefulLog = gracefulLog
}

// Tear add teardown function which calls after
func Tear(tear func() error) {
	l := instance()
	l.tears = append(l.tears, tear)
}

// Wait hold current goroutine till global context cancel.
//
// If provide wait time it will wait provided time after calling global context cancel
func Wait(waitTime ...time.Duration) {
	l := instance()

	go func() {
		signals := make(chan os.Signal)
		signal.Notify(signals, os.Interrupt, os.Kill)
		<-signals
		Cancel()
	}()

	<-l.ctx.Done()

	if l.gracefulLog != nil {
		l.gracefulLog()
	}

	tears := make([]func() error, len(l.tears))
	copy(tears, l.tears)
	slices.Reverse(tears)

	for _, tear := range tears {
		_ = try(tear)
	}

	if len(waitTime) > 0 {
		time.Sleep(waitTime[0])
	}
}

func try(fn func() error) (err error) {
	defer func() {
		if err == nil {
			err = errors.New(fmt.Sprintf("%v", recover()))
		}
	}()

	return fn()
}

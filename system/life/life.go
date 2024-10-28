package life

import (
	"context"
	"github.com/boostgo/lite/collections/list"
	"github.com/boostgo/lite/system/try"
	"os"
	"os/signal"
	"sync"
	"time"
)

var (
	_once sync.Once
	_life *life
)

type life struct {
	ctx         context.Context
	cancel      context.CancelFunc
	tears       []func() error
	gracefulLog func()
}

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

func Init() {
	instance()
}

func Context() context.Context {
	return instance().ctx
}

func Cancel() {
	instance().cancel()
}

func GracefulLog(gracefulLog func()) {
	instance().gracefulLog = gracefulLog
}

func Tear(tear func() error) {
	l := instance()
	l.tears = append(l.tears, tear)
}

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

	for _, tear := range list.Reverse(l.tears) {
		try.Must(tear)
	}

	if len(waitTime) > 0 {
		time.Sleep(waitTime[0])
	}
}

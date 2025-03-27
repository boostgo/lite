package async

import (
	"github.com/boostgo/errorx"
	"golang.org/x/sync/errgroup"
)

type Task func() error

// Go run provided function in goroutine and if panic catch recover it and convert to error
func Go(fn func()) {
	go func() {
		_ = errorx.Try(func() error {
			fn()
			return nil
		})
	}()
}

// WaitAll run all provided functions (Tasks) and wait till ending last task and return error
func WaitAll(tasks ...Task) error {
	if len(tasks) == 0 {
		return nil
	}

	wg := errgroup.Group{}
	for _, task := range tasks {
		wg.Go(func() error {
			return errorx.Try(task)
		})
	}

	return wg.Wait()
}

package storage

import (
	"context"
	"time"
)

type transactorContext struct {
	contexts []context.Context
}

func newTransactorContext(contexts ...context.Context) context.Context {
	return &transactorContext{
		contexts: contexts,
	}
}

func (trCtx *transactorContext) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

func (trCtx *transactorContext) Done() <-chan struct{} {
	return make(chan struct{}, 1)
}

func (trCtx *transactorContext) Err() error {
	for _, ctx := range trCtx.contexts {
		if err := ctx.Err(); err != nil {
			return err
		}
	}

	return nil
}

func (trCtx *transactorContext) Value(key any) any {
	for _, ctx := range trCtx.contexts {
		val := ctx.Value(key)
		if val != nil {
			return val
		}
	}

	return nil
}

package health

import (
	"context"
	"github.com/boostgo/lite/system/try"
)

type Status struct {
	CheckerName string `json:"checker_name"`
	Status      string `json:"status"`
	Error       string `json:"error,omitempty"`
}

type StatusFunc func(ctx context.Context) (Status, error)

type Checker interface {
	Status(ctx context.Context) (Status, error)
	StatusAsync(ctx context.Context) chan Status
}

type healthChecker struct {
	name       string
	getStatus  StatusFunc
	statusChan chan Status
}

func NewChecker(name string, fn StatusFunc) Checker {
	return &healthChecker{
		name:       name,
		getStatus:  fn,
		statusChan: make(chan Status, 1),
	}
}

func (checker *healthChecker) Status(ctx context.Context) (status Status, err error) {
	if err = try.Try(func() error {
		status, err = checker.getStatus(ctx)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return status, err
	}

	return Status{
		CheckerName: checker.name,
		Status:      status.Status,
		Error:       status.Error,
	}, nil
}

func (checker *healthChecker) StatusAsync(ctx context.Context) chan Status {
	go func() {
		status, err := checker.Status(ctx)
		if err != nil {
			checker.statusChan <- Status{
				CheckerName: checker.name,
				Status:      StatusUnhealthy,
				Error:       err.Error(),
			}
			return
		}

		checker.statusChan <- status
	}()
	return checker.statusChan
}

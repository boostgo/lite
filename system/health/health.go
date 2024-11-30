package health

import (
	"context"
	"github.com/boostgo/lite/app/api"
	"github.com/boostgo/lite/system/trace"
	"github.com/boostgo/lite/types/to"
	"github.com/labstack/echo/v4"
	"net/http"
	"sync"
	"time"
)

type Health struct {
	name       string
	statusPack StatusPack
	checkers   []Checker
	timeout    time.Duration
	logging    bool

	mx sync.Mutex
}

func New(name string) *Health {
	return &Health{
		name:       name,
		statusPack: StandardStatusPack(),
		checkers:   make([]Checker, 0),
	}
}

func (health *Health) Timeout(timeout time.Duration) *Health {
	health.timeout = timeout
	return health
}

func (health *Health) Logging(logging bool) *Health {
	health.logging = logging
	return health
}

func (health *Health) StatusPack(statusPack StatusPack) *Health {
	health.statusPack = statusPack
	return health
}

func (health *Health) Register(checker Checker) *Health {
	health.mx.Lock()
	defer health.mx.Unlock()

	health.checkers = append(health.checkers, checker)
	return health
}

func (health *Health) Handler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		status, statuses := health.StatusInfo()

		pretty := api.QueryParam(ctx, "pretty").Bool()
		if pretty {
			return api.Ok(ctx, map[string]any{
				"status":   status,
				"statuses": statuses,
			})
		}

		return api.SuccessRaw(ctx, http.StatusOK, to.Bytes(status))
	}
}

func (health *Health) RegisterHandler(router *echo.Echo, path string) *Health {
	router.GET(path, health.Handler(), api.Raw())
	return health
}

func (health *Health) Status() string {
	status, _ := health.StatusInfo()
	return status
}

func (health *Health) StatusInfo() (string, []Status) {
	health.mx.Lock()
	defer health.mx.Unlock()

	statuses := health.getStatuses()
	cnt := make(map[string]int)

	for _, status := range statuses {
		cnt[status.Status]++
	}

	length := len(statuses)
	switch {
	case health.statusPack.IsHealthy(cnt, length):
		return health.statusPack.Healthy, statuses
	case health.statusPack.IsUnhealthy(cnt, length):
		return health.statusPack.Unhealthy, statuses
	case health.statusPack.IsTimeout(cnt, length):
		return health.statusPack.Timeout, statuses
	}

	return health.statusPack.PartiallyUnhealthy, statuses
}

func (health *Health) getStatuses() []Status {
	ctx := trace.Set(context.Background(), trace.String())
	return newSession(ctx, health.checkers).
		Timeout(health.timeout).
		Logging(health.logging).
		Check()
}

package health

const (
	StatusHealthy            = "HEALTHY"
	StatusPartiallyUnhealthy = "PARTIALLY_UNHEALTHY"
	StatusUnhealthy          = "UNHEALTHY"
	StatusTimeout            = "TIMEOUT"
)

type StatusPack struct {
	Healthy            string
	PartiallyUnhealthy string
	Unhealthy          string
	Timeout            string
}

func StandardStatusPack() StatusPack {
	return StatusPack{
		Healthy:            StatusHealthy,
		PartiallyUnhealthy: StatusPartiallyUnhealthy,
		Unhealthy:          StatusUnhealthy,
		Timeout:            StatusTimeout,
	}
}

func (pack StatusPack) IsHealthy(cnt map[string]int, length int) bool {
	isStandard := pack.Healthy == StatusHealthy
	if isStandard {
		return cnt[StatusHealthy] == length
	}

	return cnt[StatusHealthy]+cnt[pack.Healthy] == length
}

func (pack StatusPack) IsUnhealthy(cnt map[string]int, length int) bool {
	isStandard := pack.Unhealthy == StatusUnhealthy
	if isStandard {
		return cnt[StatusUnhealthy] == length
	}

	return cnt[StatusUnhealthy]+cnt[pack.Unhealthy] == length
}

func (pack StatusPack) IsTimeout(cnt map[string]int, length int) bool {
	isStandard := pack.Timeout == StatusTimeout
	if isStandard {
		return cnt[StatusTimeout] == length
	}

	return cnt[StatusTimeout]+cnt[pack.Timeout] == length
}

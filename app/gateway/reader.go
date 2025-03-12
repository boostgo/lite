package gateway

import (
	"encoding/json"
	"os"
)

// ReadServices reads [Service] slice from .json file.
func ReadServices(path string) ([]Service, error) {
	type routeView struct {
		Method       string `json:"method" yaml:"method"`
		MatchPath    string `json:"match_path" yaml:"matchPath"`
		RedirectPath string `json:"redirect_path" yaml:"redirectPath"`
	}
	type serviceView struct {
		BaseURL string      `json:"base_url" yaml:"baseUrl"`
		Routes  []routeView `json:"routes" yaml:"routes"`
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	viewServices := make([]serviceView, 0)
	if err = json.Unmarshal(file, &viewServices); err != nil {
		return nil, err
	}

	services := make([]Service, 0, len(viewServices))
	for _, viewService := range viewServices {
		routes := make([]Route, 0, len(viewService.Routes))
		for _, viewRoute := range viewService.Routes {
			routes = append(routes, NewRoute(viewRoute.Method, viewRoute.MatchPath, viewRoute.RedirectPath))
		}

		s := NewService(viewService.BaseURL)
		s.RegisterRoute(routes...)
		services = append(services, s)
	}

	return services, nil
}

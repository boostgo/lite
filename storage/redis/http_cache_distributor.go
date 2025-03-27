package redis

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/boostgo/lite/internal/stringx"

	"github.com/boostgo/errorx"
	"github.com/boostgo/lite/api"
)

type httpCacheDistributor struct {
	client  Client
	prefix  string
	errType string
}

func NewHttpCacheDistributor(client Client, prefix string) api.HttpCacheDistributor {
	return &httpCacheDistributor{
		client:  client,
		prefix:  prefix,
		errType: "Redis HTTP Cache Distributor",
	}
}

func (distributor *httpCacheDistributor) Get(ctx context.Context, request *http.Request) (responseBody []byte, ok bool, err error) {
	defer errorx.Wrap(distributor.errType, &err, "Get")

	responseBody, err = distributor.client.GetBytes(ctx, distributor.generateKey(request))
	if err != nil {
		if errors.Is(err, errorx.ErrNotFound) {
			return nil, false, nil
		}

		return nil, false, err
	}

	return responseBody, true, nil
}

func (distributor *httpCacheDistributor) Set(ctx context.Context, request *http.Request, responseBody []byte, ttl time.Duration) (err error) {
	defer errorx.Wrap(distributor.errType, &err, "Set")
	return distributor.client.Set(ctx, distributor.generateKey(request), responseBody, ttl)
}

func (distributor *httpCacheDistributor) generateKey(request *http.Request) string {
	url := request.URL.String()
	if string(url[0]) == "/" {
		url = url[1:]
	}

	return strings.Join([]string{distributor.prefix, stringx.Replace(url, map[string]string{
		"/": "_",
		"?": "-",
		"=": "_",
		"&": "-",
	})}, "")
}

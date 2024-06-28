package config

import (
	"context"
	"github.com/boostgo/lite/types/param"
	"github.com/sethvargo/go-envconfig"
	"os"
)

func Read(export any) error {
	return ReadContext(context.Background(), export)
}

func ReadContext(ctx context.Context, export any) error {
	return envconfig.ProcessWith(ctx, export, envconfig.OsLookuper())
}

func MustReadContext(ctx context.Context, export any) {
	if err := ReadContext(ctx, export); err != nil {
		panic(err)
	}
}

func MustRead(export any) {
	MustReadContext(context.Background(), export)
}

func Get(key string) param.Param {
	return param.New(os.Getenv(key))
}

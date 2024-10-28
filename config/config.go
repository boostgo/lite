package config

import (
	"github.com/boostgo/lite/types/param"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

func Read(export any, path ...string) error {
	if len(path) == 0 {
		return cleanenv.ReadEnv(export)
	}

	for _, p := range path {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			continue
		}

		if err := cleanenv.ReadConfig(p, export); err != nil {
			return err
		}
	}

	return nil
}

func MustRead(export any, path ...string) {
	if err := Read(export, path...); err != nil {
		panic(err)
	}
}

func Get(key string) param.Param {
	return param.New(os.Getenv(key))
}

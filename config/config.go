package config

import (
	"github.com/boostgo/lite/types/param"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

func Read(export any, path ...string) error {
	var configPath string
	if len(path) > 0 {
		configPath = path[0]
	}

	return cleanenv.ReadConfig(configPath, export)
}

func MustRead(export any, path ...string) {
	if err := Read(export, path...); err != nil {
		panic(err)
	}
}

func Get(key string) param.Param {
	return param.New(os.Getenv(key))
}

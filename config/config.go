package config

import (
	"github.com/boostgo/lite/types/param"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

// Read export read config file to provided export object.
// Provided paths can contain as json/yaml file and also .env file
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

// MustRead calls Read function and if catch error throws panic
func MustRead(export any, path ...string) {
	if err := Read(export, path...); err != nil {
		panic(err)
	}
}

// Get reads OS/environment variable by provided key as param.Param object
func Get(key string) param.Param {
	return param.New(os.Getenv(key))
}

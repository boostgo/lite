package fs

import (
	"github.com/boostgo/lite/errs"
	"os"
)

func FileExist(path string) bool {
	return exist(path, existModeFile)
}

func FileWrite(path string, content []byte) error {
	if FileExist(path) {
		return errs.
			New("File already exist").
			AddContext("path", path)
	}

	return os.WriteFile(path, content, os.ModePerm)
}

func FileRead(path string) ([]byte, error) {
	if !FileExist(path) {
		return nil, errs.
			New("File not found").
			SetError(errs.ErrNotFound).
			AddContext("path", path)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func FileDelete(path string) error {
	return os.Remove(path)
}

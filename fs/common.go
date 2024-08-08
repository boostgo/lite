package fs

import (
	"os"
)

const (
	existModeBoth = iota + 1
	existModeFile
	existModeFolder
)

type Object struct {
	Name     string
	IsFolder bool
}

func newObject(name string, isFolder bool) Object {
	return Object{
		Name:     name,
		IsFolder: isFolder,
	}
}

func Exist(path string) bool {
	return exist(path, existModeBoth)
}

func Delete(path string) error {
	return remove(path, false)
}

func Rename(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

func remove(path string, withContent bool) error {
	if withContent {
		return os.RemoveAll(path)
	}

	return os.Remove(path)
}

func exist(path string, mode int) bool {
	stat, err := os.Stat(path)
	isExist := os.IsExist(err)

	switch mode {
	case existModeBoth:
		return isExist
	case existModeFile:
		return isExist && !stat.IsDir()
	case existModeFolder:
		return isExist && stat.IsDir()
	default:
		return isExist
	}
}

package fs

import (
	"github.com/boostgo/lite/errs"
	"os"
)

func FolderExist(path string) bool {
	return exist(path, existModeFolder)
}

func FolderWrite(path string) error {
	if FolderExist(path) {
		return errs.
			New("Folder already exist").
			AddContext("path", path)
	}

	return os.Mkdir(path, os.ModePerm)
}

func FolderDeleteAll(path string) error {
	return remove(path, true)
}

func FolderList(path string) ([]Object, error) {
	return folderList(path, func(_ os.DirEntry) bool {
		return true
	})
}

func FolderListFiles(path string) ([]Object, error) {
	return folderList(path, func(item os.DirEntry) bool {
		return !item.IsDir()
	})
}

func FolderListFolders(path string) ([]Object, error) {
	return folderList(path, func(item os.DirEntry) bool {
		return item.IsDir()
	})
}

func folderList(path string, condition func(item os.DirEntry) bool) ([]Object, error) {
	list, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	objects := make([]Object, 0, len(list))
	for i := 0; i < len(list); i++ {
		if !condition(list[i]) {
			continue
		}

		objects = append(objects, newObject(list[i].Name(), list[i].IsDir()))
	}

	return objects, nil
}

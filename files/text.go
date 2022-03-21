package files

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type Text struct {
	Path    string
	Content string
}

func Exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func Write(file *Text, overwrite bool) error {
	if overwrite || !Exists(file.Path) {
		data := []byte(file.Content)

		dir := filepath.Dir(file.Path)
		_ = os.MkdirAll(dir, os.ModePerm)

		return ioutil.WriteFile(file.Path, data, 0644)
	}
	return nil
}

func WriteAll(files []Text, overwrite bool) error {
	for _, file := range files {
		err := Write(&file, overwrite)
		if err != nil {
			return err
		}
	}
	return nil
}

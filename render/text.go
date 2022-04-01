package render

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type TextFile struct {
	Path    string
	Content string
}

type TextFiles []TextFile

func Exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (file *TextFile) Write(outPath string, overwrite bool) error {
	fullPath := path.Join(outPath, file.Path)
	if overwrite || !Exists(fullPath) {
		data := []byte(file.Content)

		dir := filepath.Dir(fullPath)
		_ = os.MkdirAll(dir, os.ModePerm)

		return ioutil.WriteFile(fullPath, data, 0644)
	}
	return nil
}

func (files TextFiles) WriteAll(outPath string, overwrite bool) error {
	for _, file := range files {
		err := file.Write(outPath, overwrite)
		if err != nil {
			return err
		}
	}
	return nil
}

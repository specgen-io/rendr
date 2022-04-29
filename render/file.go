package render

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

type File struct {
	Path       string
	Content    string
	Executable bool
}

type Files []File

func Exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (file *File) Write(outPath string, overwrite bool) error {
	fullPath := path.Join(outPath, file.Path)
	if overwrite || !Exists(fullPath) {
		data := []byte(file.Content)

		dir := filepath.Dir(fullPath)
		_ = os.MkdirAll(dir, os.ModePerm)

		err := ioutil.WriteFile(fullPath, data, 0644)
		if err != nil {
			return err
		}

		if file.Executable && runtime.GOOS != "windows" {
			err = os.Chmod(fullPath, 0700)
			if err != nil {
				return err
			}
		}
		return ioutil.WriteFile(fullPath, data, 0644)
	}
	return nil
}

func (files Files) WriteAll(outPath string, overwrite bool) error {
	for _, file := range files {
		err := file.Write(outPath, overwrite)
		if err != nil {
			return err
		}
	}
	return nil
}

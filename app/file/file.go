package file

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type FileDirectory struct {
	Directory string
	fileNames []string
}

func (fd FileDirectory) hasFile(search string) bool {
	for _, file := range fd.fileNames {
		if file == search {
			return true
		}
	}
	return false
}

func MakeDirectory(path string) FileDirectory {
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	files := make([]string, 0)

	for _, file := range entries {
		fmt.Println(file.Name())
		files = append(files, file.Name())
	}

	return FileDirectory{path, files}
}

func (fd FileDirectory) GetFile(file string) ([]byte, error) {
	if !fd.hasFile(file) {
		return []byte{}, errors.New("file not found")
	}

	bytes, err := os.ReadFile(fd.Directory + "/" + file)
	if err != nil {
		return bytes, err
	}
	return bytes, nil
}

func (fd FileDirectory) CreateFile(name string, data string) error {
	if fd.hasFile(name) {
		return errors.New("file already exists")
	}

	os.WriteFile(fd.Directory+"/"+name, []byte(data), os.ModeAppend)
	return nil
}

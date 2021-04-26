package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

// CreateFileIfNotExist attempts to create the full path and file if it does not exist
func CreateFileIfNotExist(filePath string) error {

	// First check to see if given filePath isn't a directory
	isDir, _ := IsDirectory(filePath)
	if isDir == true {
		return errors.New("given file path is a directory and not a path to a file")
	}

	// Check and create the base path if needed
	dir, _ := filepath.Split(filePath)
	if len(dir) > 0 {
		err := CreateDirIfNotExist(dir)
		if err != nil {
			return err
		}
	}

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		f, err := os.Create(filePath)
		defer f.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateDirIfNotExist attempts to create the directory if it does not exist
func CreateDirIfNotExist(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err == nil {
			return nil
		}
	}
	return nil
}

// IsDirectory return true if the path is a directory, false if it is a different
// type of file and and error if it doesn't exist at all
func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

func IsRegularFile(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.Mode().IsRegular(), err
}

// PathExists return true if the file/directory exists
func PathExists(path string) (bool, error) {

	_, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func IsDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// read in ONLY one file
	_, err = f.Readdir(1)

	// and if the file is EOF... well, the dir is empty.
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func NewTaskQueue() (chan<- task, <-chan task) {
	send := make(chan task)
	receive := make(chan task)
	go func() {
		queue := make([]task, 0)
		for {
			if len(queue) == 0 {
				if send == nil {
					close(receive)
					return
				}
				data, ok := <-send
				if !ok {
					close(receive)
					return
				}
				queue = append(queue, data)
			} else {
				select {
				case receive <- queue[0]:
					queue = queue[1:]
				case value, ok := <-send:
					if ok {
						queue = append(queue, value)
					} else {
						send = nil
					}
				}
			}
		}
	}()
	return send, receive
}

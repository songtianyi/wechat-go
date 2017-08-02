package rrstorage

import (
	"io/ioutil"
	"os"
	"strings"
)

// local disk storage
type LocalDiskStorage struct {
	Dir string // the directory where to save binary
}

// Create a LocalDiskStorage instance
func CreateLocalDiskStorage(dir string) StorageWrapper {
	// create dir
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		_ = os.MkdirAll(dir, os.ModeDir)
	}
	// check dir
	s := &LocalDiskStorage{
		Dir: strings.TrimSuffix(dir, "/"),
	}
	return s
}

// Do save binary
func (s *LocalDiskStorage) Save(data []byte, filename string) error {

	//open a file for writing
	file, err := os.Create(s.Dir + "/" + filename)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(data); err != nil {
		return err
	}
	return nil
}

func (s *LocalDiskStorage) Fetch(filename string) ([]byte, error) {
	b, err := ioutil.ReadFile(s.Dir + "/" + filename)
	if err != nil {
		return nil, err
	}
	return b, nil
}

package main

import (
	"bufio"
	"os"
	"path/filepath"
)

func WriteToFile(path string, b []byte) (err error) {

	f, err := os.OpenFile(filepath.Clean(path), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {

		return
	}
	defer f.Close()

	buf := bufio.NewWriter(f)
	defer buf.Flush()

	_, err = buf.Write(b)
	if err != nil {

		return
	}

	return
}

func FileExists(path string) (bool, error) {

	_, err := os.Stat(path)
	if os.IsExist(err) {

		return true, nil
	}
	if os.IsNotExist(err) {

		return false, nil
	}

	return false, err
}

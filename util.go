package main

import (
	"bufio"
	"os"
	"runtime"
)

func GetHostName() string {

	n, e := os.Hostname()
	if e != nil {

		n = "unknown"
	}
	return n
}

func GetHostOS() (os string) {

	if runtime.GOOS != "" {

		return runtime.GOOS
	}

	return "unknown"
}

func WriteToFile(file string, b []byte) (err error) {

	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0666)
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

func FileExists(file string) (bool, error) {

	_, err := os.Stat(file)
	if os.IsExist(err) {

		return true, nil
	}
	if os.IsNotExist(err) {

		return false, nil
	}

	return false, err
}

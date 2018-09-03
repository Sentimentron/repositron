package utils

import (
	"os"
)

func IsFile(path string) bool {
	if fi, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else if !fi.IsDir() {
		return true
	}
	return false
}

func IsDirectory(path string) bool {
	if fi, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else if !fi.IsDir() {
		return false
	}
	return true
}

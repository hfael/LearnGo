package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		fmt.Println("Création impossible de " + path + ". (Code 1)")
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

package main

import (
	"fmt"
	"io"
	"os"
)

func executeInit(initPath string) error {
	ok, err := exists(initPath)
	if err != nil {
		return err
	}

	if ok {
		fmt.Println("Fichier existant ! " + initPath + ".")

		fmt.Println("Lecture de " + initPath + " en cours...")

		f, err := os.Open(initPath)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		content, err := io.ReadAll(f)
		if err != nil {
			panic(err)
		}
		fmt.Println("Content: ", content)
		return nil
	}
	return nil
}

func main() {
	executeInit("init.gosql")
}

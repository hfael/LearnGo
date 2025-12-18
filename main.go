package main

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func execute(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineNo := 0

	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if err := executeLine(line); err != nil {
			return errors.New("ligne " + itoa(lineNo) + ": " + err.Error())
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func executeLine(line string) error {

	line = strings.TrimSpace(line)

	if line == "(" {
		collectSchema = true
		return nil
	}

	if line == ")" {
		if !collectSchema {
			return ErrInvalidSyntax
		}
		collectSchema = false
		return writeSchema()
	}

	if collectSchema {
		return parseColumn(line)
	}

	parts := strings.Fields(line)

	if len(parts) < 2 {
		return ErrInvalidSyntax
	}

	key := parts[0] + " " + parts[1]

	handler, ok := handlers[key]
	if !ok {
		return errors.New("instruction inconnu")
	}

	return handler(parts)
}

func main() {
	if err := execute("init.gosql"); err != nil {
		panic(err)
	}
}

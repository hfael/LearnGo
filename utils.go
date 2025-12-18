package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var ErrInvalidSyntax = errors.New("syntaxe invalide")
var ErrIncompletSyntax = errors.New("syntaxe incomplète")

type Handler func([]string) error

var handlers = map[string]Handler{
	"create database": handlerCreateDatabase,
	"create table":    handlerCreateTable,
}

var (
	currentTable    string
	currentDatabase string
	collectSchema   bool
	schemaColumns   = map[string]string{}
)

func handlerCreateDatabase(parts []string) error {

	name := strings.Trim(parts[2], `"`)
	path := strings.ToUpper(name)

	err := os.MkdirAll(path, 0775)
	if err != nil {
	}
	return err
}

func handlerCreateTable(parts []string) error {

	currentTable = strings.Trim(parts[2], `"`)
	currentDatabase = strings.Trim(parts[4], `"`)

	schemaColumns = make(map[string]string)
	collectSchema = false

	path := filepath.Join(
		strings.ToUpper(currentDatabase),
		strings.ToUpper(currentTable),
	)

	err := os.MkdirAll(path, 0775)
	if err != nil {
	}
	return err
}

func parseColumn(line string) error {

	line = strings.TrimSuffix(strings.TrimSpace(line), ",")

	parts := strings.Split(line, ":")
	if len(parts) != 2 {
		return ErrInvalidSyntax
	}

	name := strings.Trim(parts[0], `" `)
	typ := strings.TrimSpace(parts[1])

	schemaColumns[name] = typ
	return nil
}

func writeSchema() error {
	path := filepath.Join(
		strings.ToUpper(currentDatabase),
		strings.ToUpper(currentTable),
		"schema.json",
	)

	b, err := json.MarshalIndent(map[string]any{
		"table":   currentTable,
		"columns": schemaColumns,
	}, "", "  ")

	if err != nil {
		return err
	}

	err = os.WriteFile(path, b, 0644)
	if err != nil {
		return err
	}
	return nil
}

func itoa(i int) string {
	return strconv.Itoa(i)
}

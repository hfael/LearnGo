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
	"insert into":     handlerInsertInto,
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
func handlerInsertInto(parts []string) error {

	if len(parts) < 6 {
		return ErrInvalidSyntax
	}

	table := strings.Trim(parts[2], `"`)
	database := strings.Trim(parts[4], `"`)

	values := make(map[string]any)

	for _, p := range parts[6:] {
		kv := strings.SplitN(p, "=", 2)
		if len(kv) != 2 {
			return ErrInvalidSyntax
		}

		key := strings.Trim(kv[0], `"`)
		val := strings.Trim(kv[1], `"`)

		values[key] = val
	}

	return insertInto(database, table, values)
}

func insertInto(database, table string, values map[string]any) error {
	tablePath := filepath.Join(
		strings.ToUpper(database),
		strings.ToUpper(table),
	)

	schemaPath := filepath.Join(tablePath, "schema.json")

	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return err
	}

	var schema struct {
		Columns map[string]string `json:"columns"`
	}
	if err := json.Unmarshal(schemaBytes, &schema); err != nil {
		return err
	}

	primaryKey := "name"

	for col := range schema.Columns {
		if _, ok := values[col]; !ok {
			return ErrIncompletSyntax
		}
	}

	pkValue, ok := values[primaryKey]
	if !ok {
		return ErrIncompletSyntax
	}

	filename := strings.ReplaceAll(
		strings.TrimSpace(pkValue.(string)),
		" ",
		"_",
	) + ".json"

	rowPath := filepath.Join(tablePath, filename)

	out, err := json.MarshalIndent(values, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(rowPath, out, 0644)
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

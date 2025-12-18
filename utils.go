package main

import (
	"errors"
	"fmt"
	"os"
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

func handlerCreateDatabase(parts []string) error {
	if len(parts) != 3 {

		fmt.Println("3")
		return ErrInvalidSyntax
	}

	database_name := strings.Trim(parts[2], `"`)
	if database_name == "" {

		fmt.Println("4")
		return ErrInvalidSyntax
	}

	return os.Mkdir(strings.ToUpper(database_name), 0775)
}

func handlerCreateTable(parts []string) error {
	if len(parts) != 5 {
		fmt.Println("1")
		return ErrInvalidSyntax
	}
	tableName := strings.Trim(parts[2], `"`)
	databaseName := strings.Trim(parts[4], `"`)

	if tableName == "" || databaseName == "" {

		fmt.Println("2")
		return ErrInvalidSyntax
	}
	return os.Mkdir(strings.ToUpper(databaseName)+"\\"+strings.ToUpper(tableName), 0775)
}

func itoa(i int) string {
	return strconv.Itoa(i)
}

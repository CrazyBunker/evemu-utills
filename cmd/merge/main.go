package main

import (
	"fmt"
	"os"

	"game.com/m/internal/parser"
)

func main() {
	var config parser.Args
	var err error

	config, err = parser.ParseArguments(os.Args, "merge")

	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		os.Exit(1)
	}
	var base *parser.EvemuFile

	if parser.IsStdio(config.InputFile) {
		base, err = parser.ReadFromStdin()
	} else {
		base, err = parser.ParseEvemuFile(config.InputFile)
	}
	if err != nil {
		fmt.Printf("Ошибка чтения базового файла: %v\n", err)
		os.Exit(1)
	}

	add, err := parser.ParseEvemuFile(config.SecondArg)
	if err != nil {
		fmt.Printf("Ошибка чтения добавочного файла: %v\n", err)
		os.Exit(1)
	}
	// Мерж
	merged := base.Merge(add)

	if parser.IsStdio(config.OutputFile) {
		err = merged.WriteToStdout()
	} else {
		err = merged.WriteToFile(config.OutputFile)
	}

	if err != nil {
		fmt.Printf("Ошибка записи: %v\n", err)
		os.Exit(1)
	}
	if !parser.IsStdio(config.OutputFile) {
		fmt.Fprintf(os.Stderr, "Готово! Файлы объединены в %s\n", config.OutputFile)
	}
}

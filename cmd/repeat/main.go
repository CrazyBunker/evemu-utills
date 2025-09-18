package main

import (
	"fmt"
	"os"

	"game.com/m/internal/parser"
)

func main() {
	var config parser.Args
	var err error

	config, err = parser.ParseArguments(os.Args, "repeat")

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
		fmt.Printf("Ошибка чтения файла: %v\n", err)
		os.Exit(1)
	}

	// Генерация повторений
	repeated := base.GenerateRepeatedEvents(config.RepeatCount)

	if parser.IsStdio(config.OutputFile) {
		err = repeated.WriteToStdout()
	} else {
		err = repeated.WriteToFile(config.OutputFile)
	}

	if err != nil {
		fmt.Printf("Ошибка записи: %v\n", err)
		os.Exit(1)
	}
	if !parser.IsStdio(config.OutputFile) {
		fmt.Fprintf(os.Stderr, "Готово! Сгенерировано %d повторов в %s\n", config.RepeatCount, config.OutputFile)
	}
}

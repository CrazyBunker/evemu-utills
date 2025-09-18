package main

import (
	"fmt"
	"os"

	"game.com/m/internal/parser"
)

func main() {
	// Проверяем количество аргументов
	if len(os.Args) < 2 || len(os.Args) > 4 {
		fmt.Println("Использование: merge_events [базовый файл] <добавочный файл> [итоговый файл]")
		fmt.Println("  если базовый файл не указан или '-', читаем из stdin")
		fmt.Println("  если итоговый файл не указан или '-', пишем в stdout")
		os.Exit(1)
	}

	// Определяем аргументы
	var baseFile, addFile, outputFile string

	if len(os.Args) == 2 {
		baseFile = "-"
		addFile = os.Args[1]
		outputFile = "-"
	} else if len(os.Args) == 3 {
		baseFile = os.Args[1]
		addFile = os.Args[2]
		outputFile = "-"
	} else {
		baseFile = os.Args[1]
		addFile = os.Args[2]
		outputFile = os.Args[3]
	}

	// Чтение базового файла (или stdin)
	var base *parser.EvemuFile
	var err error

	if baseFile == "-" {
		base, err = parser.ReadFromStdin()
	} else {
		base, err = parser.ParseEvemuFile(baseFile)
	}

	if err != nil {
		fmt.Printf("Ошибка чтения базового файла: %v\n", err)
		os.Exit(1)
	}

	// Чтение добавочного файла (обязательно из файла)
	add, err := parser.ParseEvemuFile(addFile)
	if err != nil {
		fmt.Printf("Ошибка чтения добавочного файла: %v\n", err)
		os.Exit(1)
	}

	// Объединение файлов
	merged := base.Merge(add)

	// Запись результата
	if outputFile == "-" {
		err = merged.WriteToStdout()
	} else {
		err = merged.WriteToFile(outputFile)
	}

	if err != nil {
		fmt.Printf("Ошибка записи: %v\n", err)
		os.Exit(1)
	}
	if outputFile != "-" {
		fmt.Fprintf(os.Stderr, "Готово! Файлы объединены в %s\n", outputFile)
	}
}

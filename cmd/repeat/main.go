package main

import (
	"fmt"
	"os"
	"strconv"

	"game.com/m/internal/parser"
)

func main() {
	// Проверяем количество аргументов
	if len(os.Args) < 2 || len(os.Args) > 4 {
		fmt.Println("Использование: repeat_events [входной файл] <количество повторов> [выходной файл]")
		fmt.Println("  если входной файл не указан или '-', читаем из stdin")
		fmt.Println("  если выходной файл не указан или '-', пишем в stdout")
		os.Exit(1)
	}

	// Определяем аргументы
	var inputFile, outputFile string
	var repeatCount int
	var err error

	if len(os.Args) == 2 {
		inputFile = "-"
		repeatCount, err = strconv.Atoi(os.Args[1])
		outputFile = "-"
	} else if len(os.Args) == 3 {
		inputFile = os.Args[1]
		repeatCount, err = strconv.Atoi(os.Args[2])
		outputFile = "-"
	} else {
		inputFile = os.Args[1]
		repeatCount, err = strconv.Atoi(os.Args[2])
		outputFile = os.Args[3]
	}

	if err != nil {
		fmt.Printf("Ошибка преобразования количества повторов: %v\n", err)
		os.Exit(1)
	}

	// Чтение входного файла (или stdin)
	var file *parser.EvemuFile

	if inputFile == "-" {
		file, err = parser.ReadFromStdin()
	} else {
		file, err = parser.ParseEvemuFile(inputFile)
	}

	if err != nil {
		fmt.Printf("Ошибка чтения файла: %v\n", err)
		os.Exit(1)
	}

	// Генерация повторений
	repeated := file.GenerateRepeatedEvents(repeatCount)

	// Запись результата
	if outputFile == "-" {
		err = repeated.WriteToStdout()
	} else {
		err = repeated.WriteToFile(outputFile)
	}

	if err != nil {
		fmt.Printf("Ошибка записи: %v\n", err)
		os.Exit(1)
	}
	if outputFile != "-" {
		fmt.Fprintf(os.Stderr, "Готово! Сгенерировано %d повторов в %s\n", repeatCount, outputFile)
	}
}

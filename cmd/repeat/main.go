package main

import (
	"fmt"
	"os"
	"strconv"

	"game.com/m/internal/parser"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Использование: repeat_events <входной файл> <количество повторов> <выходной файл>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	repeatCount, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Ошибка преобразования количества повторов: %v\n", err)
		os.Exit(1)
	}
	outputFile := os.Args[3]

	file, err := parser.ParseEvemuFile(inputFile)
	if err != nil {
		fmt.Printf("Ошибка чтения файла: %v\n", err)
		os.Exit(1)
	}

	repeated := file.GenerateRepeatedEvents(repeatCount)

	if err := repeated.WriteToFile(outputFile); err != nil {
		fmt.Printf("Ошибка записи файла: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Готово! Сгенерировано %d повторов в файле %s\n", repeatCount, outputFile)
}

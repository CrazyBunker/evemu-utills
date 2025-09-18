package main

import (
	"fmt"
	"os"

	"game.com/m/internal/parser"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Использование: merge_events <базовый файл> <добавочный файл> <итоговый файл>")
		os.Exit(1)
	}

	baseFile := os.Args[1]
	addFile := os.Args[2]
	outputFile := os.Args[3]

	base, err := parser.ParseEvemuFile(baseFile)
	if err != nil {
		fmt.Printf("Ошибка чтения базового файла: %v\n", err)
		os.Exit(1)
	}

	add, err := parser.ParseEvemuFile(addFile)
	if err != nil {
		fmt.Printf("Ошибка чтения добавочного файла: %v\n", err)
		os.Exit(1)
	}

	merged := base.Merge(add)

	if err := merged.WriteToFile(outputFile); err != nil {
		fmt.Printf("Ошибка записи файла: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Готово! Файлы объединены в %s\n", outputFile)
}

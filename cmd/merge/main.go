package main

import (
	"fmt"
	"io"
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
		base, err = readFromStdin()
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
		err = writeToStdout(merged)
	} else {
		err = merged.WriteToFile(outputFile)
	}

	if err != nil {
		fmt.Printf("Ошибка записи: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Готово! Файлы объединены в %s\n", outputFile)
}

// readFromStdin читает данные из stdin и парсит их как EvemuFile
func readFromStdin() (*parser.EvemuFile, error) {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения из stdin: %v", err)
	}

	// Создаем временный файл для парсинга
	tmpfile, err := os.CreateTemp("", "stdin_input")
	if err != nil {
		return nil, fmt.Errorf("ошибка создания временного файла: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(data); err != nil {
		return nil, fmt.Errorf("ошибка записи во временный файл: %v", err)
	}

	if err := tmpfile.Close(); err != nil {
		return nil, fmt.Errorf("ошибка закрытия временного файла: %v", err)
	}

	return parser.ParseEvemuFile(tmpfile.Name())
}

// writeToStdout записывает EvemuFile в stdout
func writeToStdout(file *parser.EvemuFile) error {
	// Создаем временный файл для записи
	tmpfile, err := os.CreateTemp("", "stdout_output")
	if err != nil {
		return fmt.Errorf("ошибка создания временного файла: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if err := file.WriteToFile(tmpfile.Name()); err != nil {
		return fmt.Errorf("ошибка записи во временный файл: %v", err)
	}

	// Читаем и выводим содержимое
	data, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		return fmt.Errorf("ошибка чтения временного файла: %v", err)
	}

	fmt.Print(string(data))
	return nil
}

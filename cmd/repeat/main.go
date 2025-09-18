package main

import (
	"fmt"
	"io"
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
		file, err = readFromStdin()
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
		err = writeToStdout(repeated)
	} else {
		err = repeated.WriteToFile(outputFile)
	}

	if err != nil {
		fmt.Printf("Ошибка записи: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Готово! Сгенерировано %d повторов в %s\n", repeatCount, outputFile)
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

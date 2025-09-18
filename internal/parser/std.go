package parser

import (
	"fmt"
	"io"
	"os"
)

// readFromStdin читает данные из stdin и парсит их как EvemuFile
func ReadFromStdin() (*EvemuFile, error) {
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

	return ParseEvemuFile(tmpfile.Name())
}

// writeToStdout записывает EvemuFile в stdout
func (file *EvemuFile) WriteToStdout() error {
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

package parser

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// TestReadFromStdin тестирует чтение из stdin
func TestReadFromStdin(t *testing.T) {
	// Подготавливаем тестовые данные
	testContent := `# EVEMU 1.3
# Input device name: "Test Gamepad"
################################
E: 1.0 0001 0131 0001
E: 2.0 0001 0131 0000
`

	// Сохраняем оригинальный os.Stdin и подменяем его
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Создаем pipe для имитации stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	defer r.Close()

	os.Stdin = r

	// Записываем данные в pipe
	go func() {
		defer w.Close()
		w.Write([]byte(testContent))
	}()

	// Тестируем функцию ReadFromStdin
	result, err := ReadFromStdin()
	if err != nil {
		t.Fatalf("ReadFromStdin failed: %v", err)
	}

	// Проверяем результаты
	if len(result.Events) != 2 {
		t.Errorf("Expected 2 events, got %d", len(result.Events))
	}

	if result.Events[0].Timestamp != 1.0 || result.Events[0].Type != "0001" ||
		result.Events[0].Code != "0131" || result.Events[0].Value != "0001" {
		t.Errorf("First event mismatch: got %+v", result.Events[0])
	}
}

// TestWriteToStdout тестирует запись в stdout
func TestWriteToStdout(t *testing.T) {
	// Создаем тестовый файл
	file := &EvemuFile{
		Header: []string{"# Test header\n", "################################\n"},
		Events: []Event{
			{Timestamp: 1.234567, Type: "0001", Code: "0131", Value: "0001"},
			{Timestamp: 2.345678, Type: "0001", Code: "0131", Value: "0000"},
		},
	}

	// Перехватываем stdout
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	// Тестируем функцию WriteToStdout
	err = file.WriteToStdout()
	if err != nil {
		t.Fatalf("WriteToStdout failed: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout

	// Читаем и проверяем вывод
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	expectedLines := []string{
		"# Test header",
		"################################",
		"E: 1.234567 0001 0131 0001",
		"E: 2.345678 0001 0131 0000",
	}

	for _, line := range expectedLines {
		if !strings.Contains(output, line) {
			t.Errorf("Output missing expected line: %s", line)
		}
	}
}

// TestReadFromStdinEmpty тестирует чтение пустого stdin
func TestReadFromStdinEmpty(t *testing.T) {
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	defer r.Close()

	os.Stdin = r
	w.Close() // Закрываем сразу, пустой ввод

	result, err := ReadFromStdin()
	if err != nil {
		t.Fatalf("ReadFromStdin failed with empty input: %v", err)
	}

	if len(result.Events) != 0 {
		t.Errorf("Expected 0 events for empty input, got %d", len(result.Events))
	}
}

// TestIsStdio тестирует функцию IsStdio
func TestIsStdio(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"-", true},
		{"file.txt", false},
		{"", false},
		{"stdin", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := IsStdio(tt.path)
			if result != tt.expected {
				t.Errorf("IsStdio(%s) = %v, expected %v", tt.path, result, tt.expected)
			}
		})
	}
}

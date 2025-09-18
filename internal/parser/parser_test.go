package parser

import (
	"os"
	"path/filepath"
	"testing"
)

// TestParseEventLine тестирует разбор строки события
func TestParseEventLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected Event
		valid    bool
	}{
		{
			name:     "Valid event line",
			line:     "E: 1.234567 0003 0011 -001",
			expected: Event{Timestamp: 1.234567, Type: "0003", Code: "0011", Value: "-001"},
			valid:    true,
		},
		{
			name:     "Another valid event line",
			line:     "E: 0.583966 0001 0131 0001",
			expected: Event{Timestamp: 0.583966, Type: "0001", Code: "0131", Value: "0001"},
			valid:    true,
		},
		{
			name:  "Invalid event line - too few parts",
			line:  "E: 1.234567 0003",
			valid: false,
		},
		{
			name:  "Invalid event line - malformed timestamp",
			line:  "E: invalid 0003 0011 -001",
			valid: false,
		},
		{
			name:  "Not an event line",
			line:  "# This is a comment",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseEventLine(tt.line)

			if tt.valid {
				if result.Timestamp != tt.expected.Timestamp ||
					result.Type != tt.expected.Type ||
					result.Code != tt.expected.Code ||
					result.Value != tt.expected.Value {
					t.Errorf("parseEventLine() = %+v, expected %+v", result, tt.expected)
				}
			} else {
				// Для невалидных строк проверяем, что вернулся пустой Event
				if result.Timestamp != 0 || result.Type != "" || result.Code != "" || result.Value != "" {
					t.Errorf("parseEventLine() should return empty Event for invalid input, got %+v", result)
				}
			}
		})
	}
}

// TestParseEvemuFile тестирует разбор файла evemu
func TestParseEvemuFile(t *testing.T) {
	// Создаем временный файл для тестирования
	testContent := `# EVEMU 1.3
# Kernel: 5.15.0-112-generic
# DMI: dmi:bvnAMI:bvr7.15:bd07/02/2012:br4.6:svnHewlett-Packard:pnHPElite7500SeriesMT:pvr1.00:rvnPEGATRONCORPORATION:rn2AD5:rvr1.03:cvnHewlett-Packard:ct3:cvr:skuB5G45ES#ACB:
# Input device name: "Microsoft X-Box 360 pad"
################################
#      Waiting for events      #
################################
E: 0.000001 0003 0011 -001
E: 0.000001 0000 0000 0000
E: 0.583966 0003 0011 0000
E: 0.583966 0000 0000 0000
`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_input.txt")
	err := os.WriteFile(tmpFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Тестируем разбор файла
	result, err := ParseEvemuFile(tmpFile)
	if err != nil {
		t.Fatalf("ParseEvemuFile() failed: %v", err)
	}

	// Проверяем заголовок
	if len(result.Header) != 7 {
		t.Errorf("Expected 7 header lines, got %d", len(result.Header))
	}

	// Проверяем события
	if len(result.Events) != 4 {
		t.Errorf("Expected 4 events, got %d", len(result.Events))
	}

	// Проверяем первое событие
	firstEvent := result.Events[0]
	if firstEvent.Timestamp != 0.000001 || firstEvent.Type != "0003" || firstEvent.Code != "0011" || firstEvent.Value != "-001" {
		t.Errorf("First event mismatch: got %+v", firstEvent)
	}
}

// TestGenerateRepeatedEvents тестирует генерацию повторений событий
func TestGenerateRepeatedEvents(t *testing.T) {
	// Создаем тестовый файл с двумя событиями
	file := &EvemuFile{
		Header: []string{"# Test header\n"},
		Events: []Event{
			{Timestamp: 0.0001, Type: "0001", Code: "0131", Value: "0001"},
			{Timestamp: 0.5000, Type: "0001", Code: "0131", Value: "0000"},
		},
	}

	// Генерируем 3 повторения
	result := file.GenerateRepeatedEvents(3)

	// Проверяем количество событий
	expectedEvents := 6 // 4 исходных события × 3 повторения
	if len(result.Events) != expectedEvents {
		t.Errorf("Expected %d events, got %d", expectedEvents, len(result.Events))
	}

	// Проверяем, что заголовок сохранился
	if len(result.Header) != 1 || result.Header[0] != "# Test header\n" {
		t.Error("Header was not preserved")
	}

	// Проверяем временные метки
	expectedTimestamps := []float64{0.0, 0.49990, 0.49990, 0.99980, 0.99980, 1.49970}
	for i, event := range result.Events {
		if event.Timestamp != expectedTimestamps[i] {
			t.Errorf("Event %d: expected timestamp %.5f, got %.5f", i, expectedTimestamps[i], event.Timestamp)
		}
	}
}

// TestMerge тестирует объединение файлов
func TestMerge(t *testing.T) {
	// Создаем два тестовых файла
	baseFile := &EvemuFile{
		Header: []string{"# Base header\n"},
		Events: []Event{
			{Timestamp: 1.0, Type: "0001", Code: "0131", Value: "0001"},
			{Timestamp: 2.0, Type: "0001", Code: "0131", Value: "0000"},
		},
	}

	addFile := &EvemuFile{
		Header: []string{"# Add header\n"},
		Events: []Event{
			{Timestamp: 0.0, Type: "0003", Code: "0011", Value: "-001"},
			{Timestamp: 1.0, Type: "0003", Code: "0011", Value: "0000"},
		},
	}

	// Объединяем файлы
	result := baseFile.Merge(addFile)

	// Проверяем, что заголовок взят из базового файла
	if len(result.Header) != 1 || result.Header[0] != "# Base header\n" {
		t.Error("Header was not taken from base file")
	}

	// Проверяем количество событий
	if len(result.Events) != 4 {
		t.Errorf("Expected 4 events, got %d", len(result.Events))
	}

	// Проверяем, что временные метки скорректированы
	expectedTimestamps := []float64{1.0, 2.0, 2.0, 3.0}
	for i, event := range result.Events {
		if event.Timestamp != expectedTimestamps[i] {
			t.Errorf("Event %d: expected timestamp %.1f, got %.1f", i, expectedTimestamps[i], event.Timestamp)
		}
	}

	// Проверяем, что типы и коды событий сохранились
	if result.Events[2].Type != "0003" || result.Events[2].Code != "0011" {
		t.Error("Event types and codes were not preserved")
	}
}

// TestWriteToFile тестирует запись файла
func TestWriteToFile(t *testing.T) {
	// Создаем тестовый файл
	file := &EvemuFile{
		Header: []string{"# Test header\n", "################################\n"},
		Events: []Event{
			{Timestamp: 1.234567, Type: "0001", Code: "0131", Value: "0001"},
			{Timestamp: 2.345678, Type: "0001", Code: "0131", Value: "0000"},
		},
	}

	// Записываем во временный файл
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_output.txt")
	err := file.WriteToFile(tmpFile)
	if err != nil {
		t.Fatalf("WriteToFile() failed: %v", err)
	}

	// Читаем записанный файл и проверяем его содержимое
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expectedContent := `# Test header
################################
E: 1.234567 0001 0131 0001
E: 2.345678 0001 0131 0000
`

	if string(content) != expectedContent {
		t.Errorf("File content mismatch.\nExpected:\n%s\nGot:\n%s", expectedContent, string(content))
	}
}

// TestEmptyFiles тестирует обработку пустых файлов
func TestEmptyFiles(t *testing.T) {
	// Тестируем пустой базовый файл
	emptyBase := &EvemuFile{
		Header: []string{},
		Events: []Event{},
	}

	addFile := &EvemuFile{
		Header: []string{"# Add header\n"},
		Events: []Event{
			{Timestamp: 1.0, Type: "0001", Code: "0131", Value: "0001"},
		},
	}

	result := emptyBase.Merge(addFile)
	if len(result.Events) != 1 || result.Events[0].Timestamp != 1.0 {
		t.Error("Merge with empty base file failed")
	}

	// Тестируем пустой добавочный файл
	baseFile := &EvemuFile{
		Header: []string{"# Base header\n"},
		Events: []Event{
			{Timestamp: 1.0, Type: "0001", Code: "0131", Value: "0001"},
		},
	}

	emptyAdd := &EvemuFile{
		Header: []string{},
		Events: []Event{},
	}

	result = baseFile.Merge(emptyAdd)
	if len(result.Events) != 1 || result.Events[0].Timestamp != 1.0 {
		t.Error("Merge with empty add file failed")
	}

	// Тестируем генерацию повторений для пустого файла
	emptyFile := &EvemuFile{
		Header: []string{"# Empty header\n"},
		Events: []Event{},
	}

	result = emptyFile.GenerateRepeatedEvents(5)
	if len(result.Events) != 0 {
		t.Error("GenerateRepeatedEvents should return empty events for empty input")
	}
}

// TestParseNonExistentFile тестирует обработку несуществующего файла
func TestParseNonExistentFile(t *testing.T) {
	_, err := ParseEvemuFile("non_existent_file.txt")
	if err == nil {
		t.Error("ParseEvemuFile should return error for non-existent file")
	}
}

// TestWriteToInvalidPath тестирует запись в недопустимый путь
func TestWriteToInvalidPath(t *testing.T) {
	file := &EvemuFile{
		Header: []string{"# Test header\n"},
		Events: []Event{
			{Timestamp: 1.0, Type: "0001", Code: "0131", Value: "0001"},
		},
	}

	// Пытаемся записать в недопустимый путь
	err := file.WriteToFile("/invalid/path/test.txt")
	if err == nil {
		t.Error("WriteToFile should return error for invalid path")
	}
}

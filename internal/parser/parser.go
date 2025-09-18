package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// EvemuFile представляет файл с событиями геймпада
type EvemuFile struct {
	Header []string
	Events []Event
}

// Event представляет одно событие геймпада
type Event struct {
	Timestamp float64
	Type      string
	Code      string
	Value     string
}

// ParseEvemuFile читает и разбирает файл evemu
func ParseEvemuFile(filename string) (*EvemuFile, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла: %v", err)
	}
	defer file.Close()

	result := &EvemuFile{}
	inEventsSection := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "################################") {
			inEventsSection = true
			result.Header = append(result.Header, line+"\n")
			continue
		}

		if !inEventsSection {
			result.Header = append(result.Header, line+"\n")
		} else {
			if strings.HasPrefix(line, "E:") {
				event := parseEventLine(line)
				if event.Type != "" { // Пропускаем некорректные строки
					result.Events = append(result.Events, event)
				}
			} else {
				result.Header = append(result.Header, line+"\n")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %v", err)
	}

	return result, nil
}

// parseEventLine разбирает строку события
func parseEventLine(line string) Event {
	parts := strings.Fields(line)
	if len(parts) < 5 {
		return Event{}
	}

	timestamp, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return Event{}
	}

	return Event{
		Timestamp: timestamp,
		Type:      parts[2],
		Code:      parts[3],
		Value:     parts[4],
	}
}

// WriteToFile записывает EvemuFile в файл
func (f *EvemuFile) WriteToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Записываем заголовок
	for _, line := range f.Header {
		writer.WriteString(line)
	}

	// Записываем события
	for _, event := range f.Events {
		writer.WriteString(fmt.Sprintf("E: %.6f %s %s %s\n",
			event.Timestamp, event.Type, event.Code, event.Value))
	}

	return writer.Flush()
}

// GenerateRepeatedEvents генерирует повторения событий
func (f *EvemuFile) GenerateRepeatedEvents(repeatCount int) *EvemuFile {
	if len(f.Events) == 0 {
		return f
	}

	startTime := f.Events[0].Timestamp
	endTime := f.Events[len(f.Events)-1].Timestamp
	totalDuration := endTime - startTime

	result := &EvemuFile{
		Header: f.Header,
	}

	for i := 0; i < repeatCount; i++ {
		baseTime := float64(i) * totalDuration
		for _, event := range f.Events {
			newTime := baseTime + (event.Timestamp - startTime)
			result.Events = append(result.Events, Event{
				Timestamp: newTime,
				Type:      event.Type,
				Code:      event.Code,
				Value:     event.Value,
			})
		}
	}

	return result
}

// Merge объединяет два файла с корректировкой таймингов
func (f *EvemuFile) Merge(other *EvemuFile) *EvemuFile {
	if len(f.Events) == 0 {
		return other
	}

	if len(other.Events) == 0 {
		return f
	}

	// Находим последнее время в базовых событиях
	lastBaseTime := f.Events[len(f.Events)-1].Timestamp

	// Находим первое время в добавочных событиях
	firstAddTime := other.Events[0].Timestamp

	// Вычисляем смещение для добавочных событий
	timeOffset := lastBaseTime - firstAddTime

	// Корректируем временные метки добавочных событий
	var adjustedEvents []Event
	for _, event := range other.Events {
		adjustedEvents = append(adjustedEvents, Event{
			Timestamp: event.Timestamp + timeOffset,
			Type:      event.Type,
			Code:      event.Code,
			Value:     event.Value,
		})
	}

	// Объединяем события
	return &EvemuFile{
		Header: f.Header,
		Events: append(f.Events, adjustedEvents...),
	}
}

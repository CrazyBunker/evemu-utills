package parser

import (
	"fmt"
	"strconv"
)

type Args struct {
	InputFile   string
	SecondArg   string
	OutputFile  string
	RepeatCount int
}

func ParseArguments(args []string, utilityType string) (Args, error) {
	switch utilityType {
	case "merge":
		return parseMergeArguments(args)
	case "repeat":
		return parseRepeatArguments(args)
	default:
		return Args{InputFile: "-", SecondArg: args[1], OutputFile: "-", RepeatCount: 0}, fmt.Errorf("неизвестный тип утилиты: %s", utilityType)
	}
}

func parseMergeArguments(args []string) (Args, error) {
	if len(args) < 2 || len(args) > 4 {
		return Args{}, fmt.Errorf("использование: merge_events [базовый файл] <добавочный файл> [итоговый файл]")
	}

	if len(args) == 2 {
		return Args{InputFile: "-", SecondArg: args[1], OutputFile: "-", RepeatCount: 0}, nil
	} else if len(args) == 3 {
		return Args{InputFile: args[1], SecondArg: args[2], OutputFile: "-", RepeatCount: 0}, nil
	} else {
		return Args{InputFile: args[1], SecondArg: args[2], OutputFile: args[3], RepeatCount: 0}, nil
	}
}

func parseRepeatArguments(args []string) (Args, error) {
	if len(args) < 2 || len(args) > 4 {
		return Args{}, fmt.Errorf("использование: repeat_events [входной файл] <количество повторов> [выходной файл]")
	}

	if len(args) == 2 {
		return Args{InputFile: "-", OutputFile: "-", RepeatCount: parseRepeatCount(args[1])}, nil
	} else if len(args) == 3 {
		return Args{InputFile: args[1], OutputFile: "-", RepeatCount: parseRepeatCount(args[2])}, nil
	} else {
		return Args{InputFile: args[1], OutputFile: args[3], RepeatCount: parseRepeatCount(args[2])}, nil
	}
}

func parseRepeatCount(arg string) int {
	count, err := strconv.Atoi(arg)
	if err != nil {
		return 1
	}
	return count
}

// IsStdio проверяет, является ли путь stdin/stdout
func IsStdio(path string) bool {
	return path == "-"
}

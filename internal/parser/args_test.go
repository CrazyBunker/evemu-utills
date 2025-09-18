package parser

import (
	"testing"
)

// TestParseArgumentsMerge тестирует парсинг аргументов для merge
func TestParseArgumentsMerge(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedBase  string
		expectedAdd   string
		expectedOut   string
		expectedError bool
	}{
		{
			name:         "Valid 2 arguments",
			args:         []string{"merge", "add.txt"},
			expectedBase: "-",
			expectedAdd:  "add.txt",
			expectedOut:  "-",
		},
		{
			name:         "Valid 3 arguments",
			args:         []string{"merge", "base.txt", "add.txt"},
			expectedBase: "base.txt",
			expectedAdd:  "add.txt",
			expectedOut:  "-",
		},
		{
			name:         "Valid 4 arguments",
			args:         []string{"merge", "base.txt", "add.txt", "out.txt"},
			expectedBase: "base.txt",
			expectedAdd:  "add.txt",
			expectedOut:  "out.txt",
		},
		{
			name:          "Invalid - too few arguments",
			args:          []string{"merge"},
			expectedError: true,
		},
		{
			name:          "Invalid - too many arguments",
			args:          []string{"merge", "a", "b", "c", "d"},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := ParseArguments(tt.args, "merge")

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if config.InputFile != tt.expectedBase || config.SecondArg != tt.expectedAdd || config.OutputFile != tt.expectedOut {
				t.Errorf("ParseArguments() = (%s, %s, %s), expected (%s, %s, %s)",
					config.InputFile, config.SecondArg, config.OutputFile, tt.expectedBase, tt.expectedAdd, tt.expectedOut)
			}
		})
	}
}

// TestParseArgumentsRepeat тестирует парсинг аргументов для repeat
func TestParseArgumentsRepeat(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedInput  string
		expectedOutput string
		expectedCount  int
		expectedError  bool
	}{
		{
			name:           "Valid 2 arguments",
			args:           []string{"repeat", "5"},
			expectedInput:  "-",
			expectedOutput: "-",
			expectedCount:  5,
		},
		{
			name:           "Valid 3 arguments",
			args:           []string{"repeat", "input.txt", "3"},
			expectedInput:  "input.txt",
			expectedOutput: "-",
			expectedCount:  3,
		},
		{
			name:           "Valid 4 arguments",
			args:           []string{"repeat", "input.txt", "5", "out.txt"},
			expectedInput:  "input.txt",
			expectedOutput: "out.txt",
			expectedCount:  5,
		},
		{
			name:           "Invalid 5 arguments",
			args:           []string{"repeat", "input.txt", "5", "out.txt", "Invalid"},
			expectedInput:  "input.txt",
			expectedOutput: "out.txt",
			expectedCount:  5,
			expectedError:  true,
		},
		{
			name:           "Invalid count",
			args:           []string{"repeat", "input.txt", "invalid"},
			expectedInput:  "input.txt",
			expectedOutput: "-",
			expectedCount:  1, // default value
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := ParseArguments(tt.args, "repeat")

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if config.InputFile != tt.expectedInput || config.OutputFile != tt.expectedOutput || config.RepeatCount != tt.expectedCount {
				t.Errorf("ParseArguments() = (%s, %s, %d), expected (%s, %s, %d)",
					config.InputFile, config.OutputFile, config.RepeatCount, tt.expectedInput, tt.expectedOutput, tt.expectedCount)
			}
		})
	}
}

// TestParseArgumentsMerge тестирует парсинг аргументов для merge
func TestInvalidApplications(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedBase  string
		expectedAdd   string
		expectedOut   string
		expectedError bool
	}{
		{
			name:          "Invalid Application",
			args:          []string{"invalid", "add.txt"},
			expectedBase:  "-",
			expectedAdd:   "add.txt",
			expectedOut:   "-",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := ParseArguments(tt.args, "invalid")
			if tt.expectedError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			if config.InputFile != tt.expectedBase || config.SecondArg != tt.expectedAdd || config.OutputFile != tt.expectedOut {
				t.Errorf("ParseArguments() = (%s, %s, %s), expected (%s, %s, %s)",
					config.InputFile, config.SecondArg, config.OutputFile, tt.expectedBase, tt.expectedAdd, tt.expectedOut)
			}
		})
	}
}

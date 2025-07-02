package tools

import (
	"encoding/json"
	"testing"
)

func TestBash(t *testing.T) {
	tests := []struct {
		name        string
		input       BashInput
		expectError bool
		expectedOut string
	}{
		{
			name: "simple echo command",
			input: BashInput{
				Command: "echo 'hello world'",
			},
			expectError: false,
			expectedOut: "hello world\n",
		},
		{
			name: "pwd command",
			input: BashInput{
				Command: "pwd",
			},
			expectError: false,
		},
		{
			name: "invalid command",
			input: BashInput{
				Command: "nonexistentcommand123",
			},
			expectError: true,
		},
		{
			name: "empty command",
			input: BashInput{
				Command: "",
			},
			expectError: false,
			expectedOut: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputJSON, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("Failed to marshal input: %v", err)
			}

			result, err := Bash(inputJSON)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tt.expectedOut != "" && result != tt.expectedOut {
					t.Errorf("Expected output %q, got %q", tt.expectedOut, result)
				}
			}
		})
	}
}

func TestBashInvalidJSON(t *testing.T) {
	invalidJSON := []byte(`{"invalid": json}`)
	_, err := Bash(invalidJSON)
	if err == nil {
		t.Error("Expected error for invalid JSON input")
	}
}

func TestBashDefinition(t *testing.T) {
	if BashDefinition.Name != "bash" {
		t.Errorf("Expected name 'bash', got %q", BashDefinition.Name)
	}
	if BashDefinition.Description == "" {
		t.Error("Expected non-empty description")
	}
	if BashDefinition.Function == nil {
		t.Error("Expected non-nil function")
	}
}
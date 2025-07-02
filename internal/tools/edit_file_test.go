package tools

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestEditFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "edit_file_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name        string
		input       EditFileInput
		setupFile   func(string) error
		expectError bool
		validateFile func(string) error
	}{
		{
			name: "edit existing file",
			input: EditFileInput{
				Path:   filepath.Join(tempDir, "test.txt"),
				OldStr: "hello",
				NewStr: "hi",
			},
			setupFile: func(path string) error {
				return os.WriteFile(path, []byte("hello world"), 0644)
			},
			expectError: false,
			validateFile: func(path string) error {
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				expected := "hi world"
				if string(content) != expected {
					t.Errorf("Expected %q, got %q", expected, string(content))
				}
				return nil
			},
		},
		{
			name: "create new file",
			input: EditFileInput{
				Path:   filepath.Join(tempDir, "new.txt"),
				OldStr: "",
				NewStr: "new content",
			},
			setupFile: func(path string) error {
				return nil // No setup needed
			},
			expectError: false,
			validateFile: func(path string) error {
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				expected := "new content"
				if string(content) != expected {
					t.Errorf("Expected %q, got %q", expected, string(content))
				}
				return nil
			},
		},
		{
			name: "old_str not found",
			input: EditFileInput{
				Path:   filepath.Join(tempDir, "test2.txt"),
				OldStr: "notfound",
				NewStr: "replacement",
			},
			setupFile: func(path string) error {
				return os.WriteFile(path, []byte("hello world"), 0644)
			},
			expectError: true,
		},
		{
			name: "same old_str and new_str",
			input: EditFileInput{
				Path:   filepath.Join(tempDir, "test3.txt"),
				OldStr: "same",
				NewStr: "same",
			},
			setupFile: func(path string) error {
				return os.WriteFile(path, []byte("hello world"), 0644)
			},
			expectError: true,
		},
		{
			name: "empty path",
			input: EditFileInput{
				Path:   "",
				OldStr: "hello",
				NewStr: "hi",
			},
			setupFile: func(path string) error {
				return nil
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.setupFile(tt.input.Path); err != nil {
				t.Fatalf("Failed to setup file: %v", err)
			}

			inputJSON, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("Failed to marshal input: %v", err)
			}

			result, err := EditFile(inputJSON)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != "OK" && result == "" {
					t.Errorf("Expected 'OK' or success message, got %q", result)
				}
				if tt.validateFile != nil {
					if err := tt.validateFile(tt.input.Path); err != nil {
						t.Errorf("File validation failed: %v", err)
					}
				}
			}
		})
	}
}

func TestCreateNewFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "create_file_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test creating file in nested directory
	nestedPath := filepath.Join(tempDir, "nested", "dir", "file.txt")
	result, err := createNewFile(nestedPath, "test content")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result == "" {
		t.Error("Expected non-empty result")
	}

	// Verify file was created
	content, err := os.ReadFile(nestedPath)
	if err != nil {
		t.Errorf("Failed to read created file: %v", err)
	}
	if string(content) != "test content" {
		t.Errorf("Expected 'test content', got %q", string(content))
	}
}

func TestEditFileInvalidJSON(t *testing.T) {
	invalidJSON := []byte(`{"invalid": json}`)
	_, err := EditFile(invalidJSON)
	if err == nil {
		t.Error("Expected error for invalid JSON input")
	}
}

func TestEditFileDefinition(t *testing.T) {
	if EditFileDefinition.Name != "edit_file" {
		t.Errorf("Expected name 'edit_file', got %q", EditFileDefinition.Name)
	}
	if EditFileDefinition.Description == "" {
		t.Error("Expected non-empty description")
	}
	if EditFileDefinition.Function == nil {
		t.Error("Expected non-nil function")
	}
}
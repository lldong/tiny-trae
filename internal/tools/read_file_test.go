package tools

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestReadFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "read_file_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	testFile1 := filepath.Join(tempDir, "test1.txt")
	testContent1 := "Hello, World!\nThis is a test file."
	if err := os.WriteFile(testFile1, []byte(testContent1), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	testFile2 := filepath.Join(tempDir, "empty.txt")
	if err := os.WriteFile(testFile2, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create empty test file: %v", err)
	}

	testFile3 := filepath.Join(tempDir, "binary.bin")
	binaryContent := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE}
	if err := os.WriteFile(testFile3, binaryContent, 0644); err != nil {
		t.Fatalf("Failed to create binary test file: %v", err)
	}

	tests := []struct {
		name           string
		input          ReadFileInput
		expectError    bool
		expectedContent string
	}{
		{
			name: "read text file",
			input: ReadFileInput{
				Path: testFile1,
			},
			expectError:     false,
			expectedContent: testContent1,
		},
		{
			name: "read empty file",
			input: ReadFileInput{
				Path: testFile2,
			},
			expectError:     false,
			expectedContent: "",
		},
		{
			name: "read binary file",
			input: ReadFileInput{
				Path: testFile3,
			},
			expectError:     false,
			expectedContent: string(binaryContent),
		},
		{
			name: "read non-existent file",
			input: ReadFileInput{
				Path: filepath.Join(tempDir, "nonexistent.txt"),
			},
			expectError: true,
		},
		{
			name: "read directory instead of file",
			input: ReadFileInput{
				Path: tempDir,
			},
			expectError: true,
		},
		{
			name: "empty path",
			input: ReadFileInput{
				Path: "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputJSON, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("Failed to marshal input: %v", err)
			}

			result, err := ReadFile(inputJSON)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expectedContent {
					t.Errorf("Expected content %q, got %q", tt.expectedContent, result)
				}
			}
		})
	}
}

func TestReadFileWithRelativePath(t *testing.T) {
	// Test with relative path from current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// Create a test file in current directory
	testFile := "temp_test_file.txt"
	testContent := "temporary test content"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFile)

	input := ReadFileInput{Path: testFile}
	inputJSON, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Failed to marshal input: %v", err)
	}

	result, err := ReadFile(inputJSON)
	if err != nil {
		t.Errorf("Unexpected error reading relative path: %v", err)
	}
	if result != testContent {
		t.Errorf("Expected content %q, got %q", testContent, result)
	}

	t.Logf("Successfully read file from current directory: %s", cwd)
}

func TestReadFileInvalidJSON(t *testing.T) {
	invalidJSON := []byte(`{"invalid": json}`)
	_, err := ReadFile(invalidJSON)
	if err == nil {
		t.Error("Expected error for invalid JSON input")
	}
}

func TestReadFileDefinition(t *testing.T) {
	if ReadFileDefinition.Name != "read_file" {
		t.Errorf("Expected name 'read_file', got %q", ReadFileDefinition.Name)
	}
	if ReadFileDefinition.Description == "" {
		t.Error("Expected non-empty description")
	}
	if ReadFileDefinition.Function == nil {
		t.Error("Expected non-nil function")
	}
}

func TestReadFileLargeFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "read_large_file_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a larger test file
	testFile := filepath.Join(tempDir, "large.txt")
	largeContent := ""
	for i := 0; i < 1000; i++ {
		largeContent += "This is line " + string(rune(i)) + "\n"
	}

	if err := os.WriteFile(testFile, []byte(largeContent), 0644); err != nil {
		t.Fatalf("Failed to create large test file: %v", err)
	}

	input := ReadFileInput{Path: testFile}
	inputJSON, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Failed to marshal input: %v", err)
	}

	result, err := ReadFile(inputJSON)
	if err != nil {
		t.Errorf("Unexpected error reading large file: %v", err)
	}
	if result != largeContent {
		t.Errorf("Large file content mismatch. Expected length %d, got %d", len(largeContent), len(result))
	}
}
package tools

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRipgrep(t *testing.T) {
	// Check if ripgrep is available
	if !isRipgrepAvailable() {
		t.Skip("ripgrep (rg) is not available, skipping tests")
	}

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "ripgrep_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	testFiles := map[string]string{
		"file1.txt": "Hello World\nThis is a test file\nContains some text",
		"file2.go":  "package main\nfunc main() {\n\tfmt.Println(\"Hello World\")\n}",
		"file3.md":  "# Documentation\nThis is markdown\nHello there",
		"subdir/file4.txt": "Nested file\nHello from subdirectory\nAnother line",
	}

	for filename, content := range testFiles {
		fullPath := filepath.Join(tempDir, filename)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", fullPath, err)
		}
	}

	tests := []struct {
		name           string
		input          RipgrepInput
		expectError    bool
		expectNoMatch  bool
		expectedInOutput []string
	}{
		{
			name: "search for Hello in all files",
			input: RipgrepInput{
				Pattern: "Hello",
				Path:    tempDir,
			},
			expectError: false,
			expectedInOutput: []string{"Hello", "file1.txt", "file2.go", "file4.txt"},
		},
		{
			name: "case sensitive search",
			input: RipgrepInput{
				Pattern:       "hello",
				Path:          tempDir,
				CaseSensitive: true,
			},
			expectError:   false,
			expectNoMatch: true,
		},
		{
			name: "case insensitive search",
			input: RipgrepInput{
				Pattern:       "hello",
				Path:          tempDir,
				CaseSensitive: false,
			},
			expectError: false,
			expectedInOutput: []string{"Hello"},
		},
		{
			name: "search in specific file",
			input: RipgrepInput{
				Pattern: "package",
				Path:    filepath.Join(tempDir, "file2.go"),
			},
			expectError: false,
			expectedInOutput: []string{"package"},
		},
		{
			name: "search for non-existent pattern",
			input: RipgrepInput{
				Pattern: "nonexistentpattern123",
				Path:    tempDir,
			},
			expectError:   false,
			expectNoMatch: true,
		},
		{
			name: "search with regex pattern",
			input: RipgrepInput{
				Pattern: "[Hh]ello",
				Path:    tempDir,
			},
			expectError: false,
			expectedInOutput: []string{"Hello"},
		},
		{
			name: "search in non-existent directory",
			input: RipgrepInput{
				Pattern: "test",
				Path:    "/non/existent/path",
			},
			expectError: true,
		},
		{
			name: "search without path (current directory)",
			input: RipgrepInput{
				Pattern: "package",
			},
			expectError: false,
			// Results will vary based on current directory
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputJSON, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("Failed to marshal input: %v", err)
			}

			result, err := Ripgrep(inputJSON)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.expectNoMatch {
				if result != "No matches found." {
					t.Errorf("Expected 'No matches found.', got %q", result)
				}
				return
			}

			// For searches without specific path, just check that we got some result
			if tt.input.Path == "" {
				if result == "" {
					t.Error("Expected some result for current directory search")
				}
				return
			}

			// Check that expected strings are in the output
			for _, expected := range tt.expectedInOutput {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected %q to be in output, but got: %s", expected, result)
				}
			}
		})
	}
}

func TestRipgrepInvalidJSON(t *testing.T) {
	if !isRipgrepAvailable() {
		t.Skip("ripgrep (rg) is not available, skipping test")
	}

	invalidJSON := []byte(`{"invalid": json}`)
	_, err := Ripgrep(invalidJSON)
	if err == nil {
		t.Error("Expected error for invalid JSON input")
	}
}

func TestRipgrepDefinition(t *testing.T) {
	if RipgrepDefinition.Name != "ripgrep" {
		t.Errorf("Expected name 'ripgrep', got %q", RipgrepDefinition.Name)
	}
	if RipgrepDefinition.Description == "" {
		t.Error("Expected non-empty description")
	}
	if RipgrepDefinition.Function == nil {
		t.Error("Expected non-nil function")
	}
}

func TestRipgrepMaxCount(t *testing.T) {
	if !isRipgrepAvailable() {
		t.Skip("ripgrep (rg) is not available, skipping test")
	}

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "ripgrep_maxcount_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a file with many matches
	testFile := filepath.Join(tempDir, "many_matches.txt")
	content := ""
	for i := 0; i < 20; i++ {
		content += "test line " + string(rune(i)) + "\n"
	}
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	input := RipgrepInput{
		Pattern: "test",
		Path:    testFile,
	}
	inputJSON, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Failed to marshal input: %v", err)
	}

	result, err := Ripgrep(inputJSON)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Count the number of matches in the result
	lines := strings.Split(strings.TrimSpace(result), "\n")
	if len(lines) > 15 {
		t.Errorf("Expected at most 15 matches due to max-count, got %d", len(lines))
	}
}

// isRipgrepAvailable checks if ripgrep is available in the system
func isRipgrepAvailable() bool {
	_, err := exec.LookPath("rg")
	return err == nil
}
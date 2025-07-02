package tools

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestListFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "list_files_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files and directories
	testFiles := []string{
		"file1.txt",
		"file2.go",
		"subdir/file3.txt",
		"subdir/nested/file4.md",
	}

	for _, file := range testFiles {
		fullPath := filepath.Join(tempDir, file)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
		if err := os.WriteFile(fullPath, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", fullPath, err)
		}
	}

	tests := []struct {
		name         string
		input        ListFilesInput
		expectError  bool
		expectedFiles []string
	}{
		{
			name: "list files in temp directory",
			input: ListFilesInput{
				Path: tempDir,
			},
			expectError: false,
			expectedFiles: []string{
				"file1.txt",
				"file2.go",
				"subdir/",
				"subdir/file3.txt",
				"subdir/nested/",
				"subdir/nested/file4.md",
			},
		},
		{
			name: "list files with empty path (current directory)",
			input: ListFilesInput{
				Path: "",
			},
			expectError: false,
			// We can't predict current directory contents, so we'll just check it doesn't error
		},
		{
			name: "list files in subdirectory",
			input: ListFilesInput{
				Path: filepath.Join(tempDir, "subdir"),
			},
			expectError: false,
			expectedFiles: []string{
				"file3.txt",
				"nested/",
				"nested/file4.md",
			},
		},
		{
			name: "list files in non-existent directory",
			input: ListFilesInput{
				Path: "/non/existent/path",
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

			result, err := ListFiles(inputJSON)

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

			// Parse the JSON result
			var files []string
			if err := json.Unmarshal([]byte(result), &files); err != nil {
				t.Errorf("Failed to unmarshal result: %v", err)
				return
			}

			// For current directory test, just check that we got some result
			if tt.input.Path == "" {
				if len(files) == 0 {
					t.Error("Expected some files in current directory")
				}
				return
			}

			// Sort both slices for comparison
			sort.Strings(files)
			sort.Strings(tt.expectedFiles)

			if len(files) != len(tt.expectedFiles) {
				t.Errorf("Expected %d files, got %d. Expected: %v, Got: %v", len(tt.expectedFiles), len(files), tt.expectedFiles, files)
				return
			}

			for i, expected := range tt.expectedFiles {
				if files[i] != expected {
					t.Errorf("Expected file %q at index %d, got %q", expected, i, files[i])
				}
			}
		})
	}
}

func TestListFilesInvalidJSON(t *testing.T) {
	invalidJSON := []byte(`{"invalid": json}`)
	
	// This should panic according to the current implementation
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for invalid JSON input")
		}
	}()
	
	ListFiles(invalidJSON)
}

func TestListFilesDefinition(t *testing.T) {
	if ListFilesDefinition.Name != "list_files" {
		t.Errorf("Expected name 'list_files', got %q", ListFilesDefinition.Name)
	}
	if ListFilesDefinition.Description == "" {
		t.Error("Expected non-empty description")
	}
	if ListFilesDefinition.Function == nil {
		t.Error("Expected non-nil function")
	}
}

func TestListFilesEmptyDirectory(t *testing.T) {
	// Create an empty temporary directory
	tempDir, err := os.MkdirTemp("", "empty_dir_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	input := ListFilesInput{Path: tempDir}
	inputJSON, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Failed to marshal input: %v", err)
	}

	result, err := ListFiles(inputJSON)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var files []string
	if err := json.Unmarshal([]byte(result), &files); err != nil {
		t.Errorf("Failed to unmarshal result: %v", err)
	}

	if len(files) != 0 {
		t.Errorf("Expected empty directory to return no files, got %v", files)
	}
}
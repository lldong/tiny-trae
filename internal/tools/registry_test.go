package tools

import (
	"testing"
)

func TestGetAllTools(t *testing.T) {
	tools := GetAllTools()

	// Check that we get the expected number of tools
	expectedCount := 5
	if len(tools) != expectedCount {
		t.Errorf("Expected %d tools, got %d", expectedCount, len(tools))
	}

	// Check that all expected tools are present
	expectedTools := map[string]bool{
		"read_file":  false,
		"list_files": false,
		"edit_file":  false,
		"ripgrep":    false,
		"bash":       false,
	}

	for _, tool := range tools {
		if _, exists := expectedTools[tool.Name]; !exists {
			t.Errorf("Unexpected tool found: %s", tool.Name)
		} else {
			expectedTools[tool.Name] = true
		}

		// Validate each tool has required fields
		if tool.Name == "" {
			t.Error("Tool has empty name")
		}
		if tool.Description == "" {
			t.Errorf("Tool %s has empty description", tool.Name)
		}
		if tool.Function == nil {
			t.Errorf("Tool %s has nil function", tool.Name)
		}
		if tool.InputSchema.Type == "" {
			t.Errorf("Tool %s has empty input schema type", tool.Name)
		}
	}

	// Check that all expected tools were found
	for toolName, found := range expectedTools {
		if !found {
			t.Errorf("Expected tool %s not found in registry", toolName)
		}
	}
}

func TestGetAllToolsUniqueness(t *testing.T) {
	tools := GetAllTools()
	toolNames := make(map[string]bool)

	for _, tool := range tools {
		if toolNames[tool.Name] {
			t.Errorf("Duplicate tool name found: %s", tool.Name)
		}
		toolNames[tool.Name] = true
	}
}

func TestGetAllToolsConsistency(t *testing.T) {
	// Test that GetAllTools returns the same tools each time it's called
	tools1 := GetAllTools()
	tools2 := GetAllTools()

	if len(tools1) != len(tools2) {
		t.Errorf("GetAllTools returned different number of tools: %d vs %d", len(tools1), len(tools2))
	}

	for i, tool1 := range tools1 {
		tool2 := tools2[i]
		if tool1.Name != tool2.Name {
			t.Errorf("Tool names differ at index %d: %s vs %s", i, tool1.Name, tool2.Name)
		}
		if tool1.Description != tool2.Description {
			t.Errorf("Tool descriptions differ for %s", tool1.Name)
		}
	}
}

func TestIndividualToolDefinitions(t *testing.T) {
	// Test that individual tool definitions are properly configured
	if ReadFileDefinition.Name != "read_file" {
		t.Errorf("Expected ReadFileDefinition name 'read_file', got %q", ReadFileDefinition.Name)
	}
	if ListFilesDefinition.Name != "list_files" {
		t.Errorf("Expected ListFilesDefinition name 'list_files', got %q", ListFilesDefinition.Name)
	}
	if EditFileDefinition.Name != "edit_file" {
		t.Errorf("Expected EditFileDefinition name 'edit_file', got %q", EditFileDefinition.Name)
	}
	if RipgrepDefinition.Name != "ripgrep" {
		t.Errorf("Expected RipgrepDefinition name 'ripgrep', got %q", RipgrepDefinition.Name)
	}
	if BashDefinition.Name != "bash" {
		t.Errorf("Expected BashDefinition name 'bash', got %q", BashDefinition.Name)
	}
}
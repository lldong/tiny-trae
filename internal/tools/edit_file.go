package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"tiny-trae/internal/agent"
)

// EditFileDefinition defines the 'edit_file' tool.
var EditFileDefinition = agent.ToolDefinition{
	Name:        "edit_file",
	Description: `Make edits to a text file. Replaces 'old_str' with 'new_str' in the given file. 'old_str' and 'new_str' MUST be different from each other. If the file specified with path doesn't exist, it will be created.`,
	InputSchema: EditFileInputSchema,
	Function:    EditFile,
}

// EditFileInput defines the input schema for the 'edit_file' tool.
type EditFileInput struct {
	Path   string `json:"path" jsonschema:"description=The path to the file"`
	OldStr string `json:"old_str" jsonschema:"description=Text to search for - must match exactly and must only have one match exactly"`
	NewStr string `json:"new_str" jsonschema:"description=Text to replace old_str with"`
}

// EditFileInputSchema is the JSON schema for the 'edit_file' tool's input.
var EditFileInputSchema = agent.GenerateSchema[EditFileInput]()

// EditFile implements the 'edit_file' tool.
func EditFile(input json.RawMessage) (string, error) {
	editFileInput := EditFileInput{}
	err := json.Unmarshal(input, &editFileInput)
	if err != nil {
		return "", err
	}

	if editFileInput.Path == "" || editFileInput.OldStr == editFileInput.NewStr {
		return "", fmt.Errorf("invalid input parameters")
	}

	content, err := os.ReadFile(editFileInput.Path)
	if err != nil {
		if os.IsNotExist(err) && editFileInput.OldStr == "" {
			return createNewFile(editFileInput.Path, editFileInput.NewStr)
		}
		return "", err
	}

	oldContent := string(content)
	newContent := strings.Replace(oldContent, editFileInput.OldStr, editFileInput.NewStr, -1)

	if oldContent == newContent && editFileInput.OldStr != "" {
		return "", fmt.Errorf("old_str not found in file")
	}

	err = os.WriteFile(editFileInput.Path, []byte(newContent), 0644)
	if err != nil {
		return "", err
	}

	return "OK", nil
}

// createNewFile creates a new file with the given content.
func createNewFile(filePath, content string) (string, error) {
	dir := path.Dir(filePath)
	if dir != "." {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}
	}

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}

	return fmt.Sprintf("Successfully created file %s", filePath), nil
}

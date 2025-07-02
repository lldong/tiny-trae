package tools

import (
	"encoding/json"
	"os"

	"tiny-trae/internal/agent"
)

// ReadFileDefinition defines the 'read_file' tool.
var ReadFileDefinition = agent.ToolDefinition{
	Name:        "read_file",
	Description: "Read the contents of a given relative file path. Use this when you want to see what's inside a file. Do not use this with directory names.",
	InputSchema: ReadFileInputSchema,
	Function:    ReadFile,
}

// ReadFileInput defines the input schema for the 'read_file' tool.
type ReadFileInput struct {
	Path string `json:"path" jsonschema:"description=The relative path of a file in the working directory"`
}

// ReadFileInputSchema is the JSON schema for the 'read_file' tool's input.
var ReadFileInputSchema = agent.GenerateSchema[ReadFileInput]()

// ReadFile implements the 'read_file' tool.
func ReadFile(input json.RawMessage) (string, error) {
	readFileInput := ReadFileInput{}
	err := json.Unmarshal(input, &readFileInput)
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(readFileInput.Path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

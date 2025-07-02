package tools

import (
	"encoding/json"
	"os"
	"path/filepath"

	"tiny-trae/internal/agent"
)

// ListFilesDefinition defines the 'list_files' tool.
var ListFilesDefinition = agent.ToolDefinition{
	Name:        "list_files",
	Description: "List files and directories at a given path. If no path is provided, lists files in the current directory.",
	InputSchema: ListFilesInputSchema,
	Function:    ListFiles,
}

// ListFilesInput defines the input schema for the 'list_files' tool.
type ListFilesInput struct {
	Path string `json:"path,omitempty" jsonschema:"description=Optional relative path to list files from. Defaults to current directory if not provided."`
}

// ListFilesInputSchema is the JSON schema for the 'list_files' tool's input.
var ListFilesInputSchema = agent.GenerateSchema[ListFilesInput]()

// ListFiles implements the 'list_files' tool.
func ListFiles(input json.RawMessage) (string, error) {
	listFilesInput := ListFilesInput{}
	err := json.Unmarshal(input, &listFilesInput)
	if err != nil {
		panic(err)
	}

	dir := "."
	if listFilesInput.Path != "" {
		dir = listFilesInput.Path
	}

	var files []string
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		if relPath != "." {
			if info.IsDir() {
				files = append(files, relPath+"/")
			} else {
				files = append(files, relPath)
			}
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	result, err := json.Marshal(files)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

package tools

import "tiny-trae/internal/agent"

// GetAllTools returns all available tool definitions.
func GetAllTools() []agent.ToolDefinition {
	return []agent.ToolDefinition{
		ReadFileDefinition,
		ListFilesDefinition,
		EditFileDefinition,
		RipgrepDefinition,
		BashDefinition,
	}
}

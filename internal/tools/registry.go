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

// GetMinimalTools returns a minimal set of tools for basic tasks.
func GetMinimalTools() []agent.ToolDefinition {
	return []agent.ToolDefinition{
		ReadFileDefinition,
		ListFilesDefinition,
		EditFileDefinition,
	}
}

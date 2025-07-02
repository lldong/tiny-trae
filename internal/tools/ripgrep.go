package tools

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"tiny-trae/internal/agent"
)

// RipgrepDefinition defines the 'ripgrep' tool.
var RipgrepDefinition = agent.ToolDefinition{
	Name: "ripgrep",
	Description: `Search for exact text patterns in files using ripgrep, a fast keyword search tool.

WHEN TO USE THIS TOOL:
- When you need to find exact text matches like variable names, function calls, or specific strings
- When you know the precise pattern you're looking for (including regex patterns)
- When you want to quickly locate all occurrences of a specific term across multiple files
- When you need to search for code patterns with exact syntax

WHEN NOT TO USE THIS TOOL:
- For semantic or conceptual searches (e.g., "how does authentication work")
- For finding code that implements a certain functionality without knowing the exact terms
- When you already have read the entire file

RESULT INTERPRETATION:
- Results show the file path, line number, and matching line content
- Results are grouped by file, with up to 15 matches per file`,
	InputSchema: RipgrepInputSchema,
	Function:    Ripgrep,
}

// RipgrepInput defines the input schema for the 'ripgrep' tool.
type RipgrepInput struct {
	Pattern       string `json:"pattern" jsonschema_description:"The pattern to search for"`
	Path          string `json:"path,omitempty" jsonschema_description:"The file or directory path to search in"`
	CaseSensitive bool   `json:"caseSensitive,omitempty" jsonschema_description:"Whether to search case-sensitively"`
}

// RipgrepInputSchema is the JSON schema for the 'ripgrep' tool's input.
var RipgrepInputSchema = agent.GenerateSchema[RipgrepInput]()

// Ripgrep implements the 'ripgrep' tool.
func Ripgrep(input json.RawMessage) (string, error) {
	ripgrepInput := RipgrepInput{}
	err := json.Unmarshal(input, &ripgrepInput)
	if err != nil {
		return "", err
	}

	args := []string{"--line-number"}

	if !ripgrepInput.CaseSensitive {
		args = append(args, "-i")
	}

	args = append(args, "--max-count", "15")
	args = append(args, ripgrepInput.Pattern)

	if ripgrepInput.Path != "" {
		args = append(args, ripgrepInput.Path)
	}

	cmd := exec.Command("rg", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		// Exit code 1 in ripgrep means "no matches found", which isn't an error for us
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return "No matches found.", nil
		}
		return "", fmt.Errorf("ripgrep error: %v - %s", err, string(output))
	}

	return string(output), nil
}

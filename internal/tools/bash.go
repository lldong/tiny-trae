package tools

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"tiny-trae/internal/agent"
)

// BashDefinition defines the 'bash' tool.
var BashDefinition = agent.ToolDefinition{
	Name:        "bash",
	Description: "Execute a bash command.",
	InputSchema: BashInputSchema,
	Function:    Bash,
}

// BashInput defines the input schema for the 'bash' tool.
type BashInput struct {
	Command string `json:"command" jsonschema:"description=The command to execute"`
}

// BashInputSchema is the JSON schema for the 'bash' tool's input.
var BashInputSchema = agent.GenerateSchema[BashInput]()

// Bash implements the 'bash' tool.
func Bash(input json.RawMessage) (string, error) {
	bashInput := BashInput{}
	err := json.Unmarshal(input, &bashInput)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("bash", "-c", bashInput.Command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command execution error: %v - %s", err, string(output))
	}

	return string(output), nil
}

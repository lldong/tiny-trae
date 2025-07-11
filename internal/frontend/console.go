package frontend

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"tiny-trae/internal/agent"
)

// ConsoleFrontend implements the Frontend interface for console-based interaction
type ConsoleFrontend struct {
	scanner     *bufio.Scanner
	interactive bool
}

// NewConsoleFrontend creates a new console frontend
func NewConsoleFrontend(interactive bool) *ConsoleFrontend {
	return &ConsoleFrontend{
		scanner:     bufio.NewScanner(os.Stdin),
		interactive: interactive,
	}
}

// SendMessage sends a message to the console for display
func (c *ConsoleFrontend) SendMessage(msg agent.Message) {
	switch msg.Type {
	case agent.MessageTypeUserInput:
		// Don't display user input here as it's already shown in GetUserInput
		// This prevents duplicate output
	case agent.MessageTypeAssistant:
		if c.interactive {
			fmt.Printf("Trae: %s\n", msg.Content)
		} else {
			fmt.Printf("%s\n", msg.Content)
		}
	case agent.MessageTypeToolCall:
		if c.interactive {
			var toolData agent.ToolCallData
			if err := json.Unmarshal(msg.Data, &toolData); err == nil {
				fmt.Printf("Tool: %s(%s)\n", toolData.ToolName, toolData.Input)
			}
		}
	case agent.MessageTypeToolResult:
		// Display tool results in interactive mode
		var toolResult agent.ToolResultData
		if err := json.Unmarshal(msg.Data, &toolResult); err == nil {
			if c.interactive {
				if toolResult.IsError {
					fmt.Printf("Tool Error (%s): %s\n", toolResult.ToolName, toolResult.Result)
				} else {
					fmt.Printf("Tool Result (%s): %s\n", toolResult.ToolName, toolResult.Result)
				}
			}
		}
	case agent.MessageTypeError:
		fmt.Printf("Error: %s\n", msg.Content)
	case agent.MessageTypeSystemInfo:
		fmt.Printf("%s\n", msg.Content)
	}
}

// GetUserInput requests user input from the console
func (c *ConsoleFrontend) GetUserInput() (string, bool) {
	if !c.interactive {
		return "", false
	}

	fmt.Print("You: ")
	if !c.scanner.Scan() {
		if err := c.scanner.Err(); err != nil {
			fmt.Printf("Error reading input: %v\n", err)
		}
		return "", false
	}
	return c.scanner.Text(), true
}

// IsInteractive returns whether the console frontend is in interactive mode
func (c *ConsoleFrontend) IsInteractive() bool {
	return c.interactive
}

// Close closes the console frontend
func (c *ConsoleFrontend) Close() {
	// Nothing special to do for console frontend
}

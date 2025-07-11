package agent

import "encoding/json"

// MessageType represents the type of message being sent to the frontend
type MessageType string

const (
	MessageTypeUserInput    MessageType = "user_input"
	MessageTypeAssistant    MessageType = "assistant"
	MessageTypeToolCall     MessageType = "tool_call"
	MessageTypeToolResult   MessageType = "tool_result"
	MessageTypeError        MessageType = "error"
	MessageTypeSystemInfo   MessageType = "system_info"
)

// Message represents a message sent from the agent core to the frontend
type Message struct {
	Type    MessageType     `json:"type"`
	Content string          `json:"content"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// ToolCallData represents additional data for tool call messages
type ToolCallData struct {
	ToolName string          `json:"tool_name"`
	ToolID   string          `json:"tool_id"`
	Input    json.RawMessage `json:"input"`
}

// ToolResultData represents additional data for tool result messages
type ToolResultData struct {
	ToolName string `json:"tool_name"`
	ToolID   string `json:"tool_id"`
	Result   string `json:"result"`
	IsError  bool   `json:"is_error"`
}

// Frontend represents the interface that any frontend implementation must satisfy
type Frontend interface {
	// SendMessage sends a message to the frontend for display
	SendMessage(msg Message)
	// GetUserInput requests user input from the frontend
	GetUserInput() (string, bool)
	// IsInteractive returns whether the frontend is in interactive mode
	IsInteractive() bool
	// Close closes the frontend
	Close()
}
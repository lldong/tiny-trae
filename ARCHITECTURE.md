# Tiny Trae Architecture

## Overview

Tiny Trae has been refactored to separate the core agent logic from the frontend interface, making it easy to add different frontend implementations.

## Architecture Components

### Core Components

1. **Agent Core** (`internal/agent/agent.go`)
   - Contains the main AI agent logic
   - Runs in a separate goroutine
   - Communicates with frontend through message passing
   - Handles conversation management and tool execution

2. **Message System** (`internal/agent/message.go`)
   - Defines message types for communication between core and frontend
   - Provides structured data for different message types (user input, assistant response, tool calls, etc.)

3. **Frontend Interface** (`internal/agent/message.go`)
   - Defines the `Frontend` interface that all frontend implementations must satisfy
   - Provides methods for sending messages and getting user input

### Current Frontend Implementation

4. **TUI Frontend** (`internal/frontend/tui.go`)
- Implements the `Frontend` interface for terminal user interface interaction using bubbletea
   - Handles user input from stdin and displays messages to stdout
   - Supports both interactive and non-interactive modes

## Message Types

The system uses the following message types for communication:

- `MessageTypeUserInput`: User input messages
- `MessageTypeAssistant`: AI assistant responses
- `MessageTypeToolCall`: Tool execution notifications
- `MessageTypeToolResult`: Tool execution results
- `MessageTypeError`: Error messages
- `MessageTypeSystemInfo`: System information messages

## How to Add a New Frontend

To add a new frontend (e.g., web interface, GUI, API server), follow these steps:

1. **Create a new frontend package** (e.g., `internal/frontend/web.go`)

2. **Implement the Frontend interface**:
   ```go
   type YourFrontend struct {
       // Your frontend-specific fields
   }

   func (f *YourFrontend) SendMessage(msg agent.Message) {
       // Handle different message types and display them appropriately
       switch msg.Type {
       case agent.MessageTypeUserInput:
           // Handle user input display
       case agent.MessageTypeAssistant:
           // Handle assistant response display
       case agent.MessageTypeToolCall:
           // Handle tool call notifications
       // ... handle other message types
       }
   }

   func (f *YourFrontend) GetUserInput() (string, bool) {
       // Get user input from your frontend
       // Return the input string and a boolean indicating success
   }

   func (f *YourFrontend) Close() {
       // Clean up resources
   }
   ```

3. **Update main.go** to use your new frontend:
   ```go
   // Replace the TUI frontend creation with your frontend
   yourFrontend := frontend.NewYourFrontend(/* your parameters */)
   defer yourFrontend.Close()
   
   // Create agent with your frontend
   agentInstance := agent.NewAgent(client, agentProfile, yourFrontend)
   ```

## Benefits of This Architecture

1. **Separation of Concerns**: Core logic is separated from UI concerns
2. **Concurrent Execution**: Agent core runs in its own goroutine
3. **Extensibility**: Easy to add new frontend implementations
4. **Testability**: Core logic can be tested independently of UI
5. **Flexibility**: Different frontends can handle messages differently based on their capabilities

## Example Frontend Implementations

Possible frontend implementations you could add:

- **Web Frontend**: HTTP server with WebSocket for real-time communication
- **GUI Frontend**: Desktop application using frameworks like Fyne or Qt
- **API Frontend**: REST API server for integration with other applications
- **Chat Bot Frontend**: Integration with Discord, Slack, or other chat platforms
- **Mobile Frontend**: Mobile app interface

## Running the Application

```bash
# Interactive mode with default profile
./tiny-trae

# Non-interactive mode with a single prompt
./tiny-trae -p "Hello, how are you?"

# Use a specific profile
./tiny-trae -profile minimal

# List available profiles
./tiny-trae --list-profiles
```
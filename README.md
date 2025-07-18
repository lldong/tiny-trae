# Tiny Trae

A minimal AI coding agent powered by Anthropic's Claude with a modular frontend architecture.

🌐 **Landing Page**: https://lldong.github.io/tiny-trae

## Architecture

Tiny Trae has been designed with a clean separation between the core agent logic and the frontend interface. This allows for easy extension with different types of user interfaces.

### Core Components

- **Agent Core**: Handles AI conversation logic and tool execution in a separate goroutine
- **Frontend Interface**: Defines how different UIs can interact with the agent
- **Message System**: Structured communication between core and frontend
- **TUI Frontend**: Terminal user interface implementation using bubbletea

### Frontend

- **TUI**: Terminal user interface with rich interface using bubbletea

For detailed architecture information, see [ARCHITECTURE.md](ARCHITECTURE.md).

This project is a simple AI coding agent implemented in Go. It uses the Anthropic API to interact with a large language model (Claude) to help with software engineering tasks. The agent can execute a predefined set of tools based on the model's response.

## Features

- **Interactive Chat:** Chat with the agent from your terminal.
- **Non-interactive mode:** Provide input directly from the command line.
- **Tool Execution:** The agent can execute the following tools:
    - `read_file`: Read the contents of a file.
    - `list_files`: List files and directories.
    - `edit_file`: Modify files by searching and replacing text.
    - `ripgrep`: Search for text patterns within files.
    - `bash`: Execute shell commands.
- **Extensible:** Easily add new tools to the agent.

## Prerequisites

- Go 1.x
- An [Anthropic API key](https://console.anthropic.com/dashboard)
- **ripgrep**: This tool is used by the `ripgrep` command. You can install it by following the instructions in the [ripgrep repository](https://github.com/BurntSushi/ripgrep#installation). For example, on macOS you can use Homebrew: `brew install ripgrep`

## Getting Started

1.  **Clone the repository:**
    ```bash
    git clone <repository-url>
    cd tiny-trae
    ```

2.  **Install dependencies:**
    ```bash
    go mod tidy
    ```

3.  **Set up your Anthropic API key:**
    ```bash
    export ANTHROPIC_API_KEY="your-api-key"
    ```
    You can also set `ANTHROPIC_BASE_URL` if you are using a proxy.

4.  **Run the agent:**
    ```bash
    go run main.go
    ```
    or
    ```bash
    go build -o tiny-trae
    ```

## Usage

### Interactive Mode

To run the agent in interactive mode, simply run the executable:

```bash
./tiny-trae
```

The agent will prompt you for input.

### Non-interactive Mode

To run the agent in non-interactive mode, use the `-p` flag to provide a prompt:

```bash
./tiny-trae -p "your prompt here"
```

The agent will process the prompt and exit.

## Using with OpenRouter

You can use this agent with [OpenRouter](https://openrouter.ai/) by using `anthropic-proxy`.

1.  **Start the proxy:**
    ```bash
    OPENROUTER_API_KEY=your-api-key COMPLETION_MODEL="anthropic/claude-sonnet-4" npx anthropic-proxy
    ```

2.  **Run the agent:**
    In a separate terminal, run the following command:
    ```bash
    ANTHROPIC_BASE_URL=http://0.0.0.0:3000 ./tiny-trae
    ```

## How it works

The agent starts a conversation with the user. The user's message is sent to the Anthropic API, and the model can either respond with text or a request to use a tool. If it's a tool-use request, the agent executes the tool and sends the result back to the model. This loop continues until the user exits the program.

## Tools

The agent currently supports the following tools:

-   **`read_file`**: Reads the entire content of a specified file.
-   **`list_files`**: Lists all files and directories within a given path.
-   **`edit_file`**: Edits a file by replacing a specified string with a new one.
-   **`ripgrep`**: Searches for a pattern in files using `rg`.
-   **`bash`**: Executes a given command in a bash shell.

You can extend the agent by adding new `ToolDefinition` structs and including them in the `tools` slice in the `main` function.

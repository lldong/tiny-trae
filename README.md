# Tiny Trae

This project is a simple AI coding agent implemented in Go. It uses the Anthropic API to interact with a large language model (Claude) to help with software engineering tasks. The agent can execute a predefined set of tools based on the model's response.

## Features

- **Interactive Chat:** Chat with the agent from your terminal.
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
    ./tiny-trae
    ```

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

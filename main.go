package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/invopop/jsonschema"
)

// SYSTEM_PROMPT is the default prompt that defines the agent's persona.
const SYSTEM_PROMPT = `You are a powerful AI coding agent. You help the user with software engineering tasks.
`

// Agent struct represents the core of the AI agent.
type Agent struct {
	client         *anthropic.Client
	getUserMessage func() (string, bool)
	tools          []ToolDefinition
	interactive    bool
}

// ToolDefinition struct defines a tool that the agent can use.
type ToolDefinition struct {
	Name        string                         `json:"name"`
	Description string                         `json:"description"`
	InputSchema anthropic.ToolInputSchemaParam `json:"input_schema"`
	Function    func(input json.RawMessage) (string, error)
}

// NewAgent creates a new Agent instance.
func NewAgent(
	client *anthropic.Client,
	getUserMessage func() (string, bool),
	tools []ToolDefinition,
	interactive bool,
) *Agent {
	return &Agent{
		client:         client,
		getUserMessage: getUserMessage,
		tools:          tools,
		interactive:    interactive,
	}
}

// Run starts the agent's main loop.
// It continuously prompts the user for input, sends it to the Anthropic API,
// and processes the model's response, which may include text or tool use requests.
// The loop terminates when the user signals the end of input (e.g., by pressing CTRL+C).
// In non-interactive mode, it takes an initial message, gets the model's response, and exits.
func (a *Agent) Run(ctx context.Context, initialMessage string) error {
	conversation := []anthropic.MessageParam{}

	if initialMessage != "" {
		userMessage := anthropic.NewUserMessage(anthropic.NewTextBlock(initialMessage))
		conversation = append(conversation, userMessage)
	} else {
		fmt.Printf("Chat with Tiny Trae (use CTRL+C to exit)\n")
	}

	readUserInput := initialMessage == ""
	for {
		if readUserInput {
			userInput, ok := a.getUserMessage()
			if !ok {
				break
			}

			userMessage := anthropic.NewUserMessage(anthropic.NewTextBlock(userInput))
			conversation = append(conversation, userMessage)
		}

		message, err := a.runInference(ctx, conversation)
		if err != nil {
			return err
		}
		conversation = append(conversation, message.ToParam())

		hasToolUse := false
		for _, content := range message.Content {
			if content.Type == "tool_use" {
				hasToolUse = true
				break
			}
		}

		toolResults := []anthropic.ContentBlockParamUnion{}
		for _, content := range message.Content {
			switch content.Type {
			case "text":
				if a.interactive {
					fmt.Printf("Trae: %s\n", content.Text)
				} else if !hasToolUse {
					fmt.Printf("%s\n", content.Text)
				}
			case "tool_use":
				result := a.executeTool(content.ID, content.Name, content.Input)
				toolResults = append(toolResults, result)
			}
		}

		if len(toolResults) == 0 {
			if !a.interactive {
				return nil
			}
			readUserInput = true
			continue
		}

		readUserInput = false
		conversation = append(conversation, anthropic.NewUserMessage(toolResults...))
	}

	return nil
}

// runInference sends the conversation to the Anthropic API and gets the model's response.
// It constructs a list of tools available for the model to use and includes them in the API request.
// The function returns the model's response message or an error if the API call fails.
func (a *Agent) runInference(ctx context.Context, conversation []anthropic.MessageParam) (*anthropic.Message, error) {
	anthropicTools := []anthropic.ToolUnionParam{}
	for _, tool := range a.tools {
		anthropicTools = append(anthropicTools, anthropic.ToolUnionParam{
			OfTool: &anthropic.ToolParam{
				Name:        tool.Name,
				Description: anthropic.String(tool.Description),
				InputSchema: tool.InputSchema,
			},
		})
	}

	message, err := a.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_0,
		MaxTokens: int64(1024),
		Messages:  conversation,
		Tools:     anthropicTools,
		System:    []anthropic.TextBlockParam{{Text: SYSTEM_PROMPT}},
	})

	return message, err
}

// executeTool executes a tool with the given name and input.
// It finds the corresponding tool definition, calls its associated function with the provided input,
// and returns the result as a tool result block. If the tool is not found or an error occurs
// during execution, it returns an error message in the tool result block.
func (a *Agent) executeTool(id, name string, input json.RawMessage) anthropic.ContentBlockParamUnion {
	var toolDef ToolDefinition
	var found bool
	for _, tool := range a.tools {
		if tool.Name == name {
			toolDef = tool
			found = true
			break
		}
	}
	if !found {
		return anthropic.NewToolResultBlock(id, "tool not found", true)
	}

	if a.interactive {
		fmt.Printf("Tool: %s(%s)\n", name, input)
	}

	response, err := toolDef.Function(input)
	if err != nil {
		return anthropic.NewToolResultBlock(id, err.Error(), true)
	}
	return anthropic.NewToolResultBlock(id, response, false)
}

// ReadFileDefinition defines the 'read_file' tool.
var ReadFileDefinition = ToolDefinition{
	Name:        "read_file",
	Description: "Read the contents of a given relative file path. Use this when you want to see what's inside a file. Do not use this with directory names.",
	InputSchema: ReadFileInputSchema,
	Function:    ReadFile,
}

// ReadFileInput defines the input schema for the 'read_file' tool.
type ReadFileInput struct {
	Path string `json:"path" jsonschema:"description=The relative path of a file in the working directory"`
}

// ReadFileInputSchema is the JSON schema for the 'read_file' tool's input.
var ReadFileInputSchema = GenerateSchema[ReadFileInput]()

// ReadFile implements the 'read_file' tool.
func ReadFile(input json.RawMessage) (string, error) {
	readFileInput := ReadFileInput{}
	err := json.Unmarshal(input, &readFileInput)
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(readFileInput.Path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// ListFilesDefinition defines the 'list_files' tool.
var ListFilesDefinition = ToolDefinition{
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
var ListFilesInputSchema = GenerateSchema[ListFilesInput]()

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

// EditFileDefinition defines the 'edit_file' tool.
var EditFileDefinition = ToolDefinition{
	Name:        "edit_file",
	Description: `Make edits to a text file. Replaces 'old_str' with 'new_str' in the given file. 'old_str' and 'new_str' MUST be different from each other. If the file specified with path doesn't exist, it will be created.`,
	InputSchema: EditFileInputSchema,
	Function:    EditFile,
}

// EditFileInput defines the input schema for the 'edit_file' tool.
type EditFileInput struct {
	Path   string `json:"path" jsonschema:"description=The path to the file"`
	OldStr string `json:"old_str" jsonschema:"description=Text to search for - must match exactly and must only have one match exactly"`
	NewStr string `json:"new_str" jsonschema:"description=Text to replace old_str with"`
}

// EditFile implements the 'edit_file' tool.
func EditFile(input json.RawMessage) (string, error) {
	editFileInput := EditFileInput{}
	err := json.Unmarshal(input, &editFileInput)
	if err != nil {
		return "", err
	}

	if editFileInput.Path == "" || editFileInput.OldStr == editFileInput.NewStr {
		return "", fmt.Errorf("invalid input parameters")
	}

	content, err := os.ReadFile(editFileInput.Path)
	if err != nil {
		if os.IsNotExist(err) && editFileInput.OldStr == "" {
			return createNewFile(editFileInput.Path, editFileInput.NewStr)
		}
		return "", err
	}

	oldContent := string(content)
	newContent := strings.Replace(oldContent, editFileInput.OldStr, editFileInput.NewStr, -1)

	if oldContent == newContent && editFileInput.OldStr != "" {
		return "", fmt.Errorf("old_str not found in file")
	}

	err = os.WriteFile(editFileInput.Path, []byte(newContent), 0644)
	if err != nil {
		return "", err
	}

	return "OK", nil
}

// createNewFile creates a new file with the given content.
func createNewFile(filePath, content string) (string, error) {
	dir := path.Dir(filePath)
	if dir != "." {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}
	}

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}

	return fmt.Sprintf("Successfully created file %s", filePath), nil
}

// EditFileInputSchema is the JSON schema for the 'edit_file' tool's input.
var EditFileInputSchema = GenerateSchema[EditFileInput]()

// GenerateSchema generates a JSON schema for a given type.
func GenerateSchema[T any]() anthropic.ToolInputSchemaParam {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}

	var v T
	schema := reflector.Reflect(v)

	return anthropic.ToolInputSchemaParam{
		Type:       "object",
		Properties: schema.Properties,
	}
}

// RipgrepDefinition defines the 'ripgrep' tool.
var RipgrepDefinition = ToolDefinition{
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
var RipgrepInputSchema = GenerateSchema[RipgrepInput]()

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

// BashDefinition defines the 'bash' tool.
var BashDefinition = ToolDefinition{
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
var BashInputSchema = GenerateSchema[BashInput]()

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

// main is the entry point of the application.
// It initializes the Anthropic client, sets up the available tools,
// creates a new agent, and starts its execution.
// It supports both interactive and non-interactive modes.
// Any errors that occur during the agent's run are printed to the console.
func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println()
		os.Exit(0)
	}()

	var options []option.RequestOption
	if apiKey := os.Getenv("ANTHROPIC_API_KEY"); apiKey != "" {
		options = append(options, option.WithAPIKey(apiKey))
	}
	if baseURL := os.Getenv("ANTHROPIC_BASE_URL"); baseURL != "" {
		options = append(options, option.WithBaseURL(baseURL))
	}
	client := anthropic.NewClient(options...)

	prompt := flag.String("p", "", "Accept a string as user input")
	flag.Parse()

	var getUserMessage func() (string, bool)
	var initialMessage string

	if *prompt != "" {
		initialMessage = *prompt
		getUserMessage = func() (string, bool) {
			return "", false
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		getUserMessage = func() (string, bool) {
			fmt.Print("You: ")
			if !scanner.Scan() {
				return "", false
			}
			return scanner.Text(), true
		}
	}

	tools := []ToolDefinition{
		ReadFileDefinition,
		ListFilesDefinition,
		EditFileDefinition,
		RipgrepDefinition,
		BashDefinition,
	}
	agent := NewAgent(&client, getUserMessage, tools, *prompt == "")
	err := agent.Run(context.TODO(), initialMessage)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}

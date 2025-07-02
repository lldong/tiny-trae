package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/invopop/jsonschema"
)

// ToolDefinition struct defines a tool that the agent can use.
type ToolDefinition struct {
	Name        string                         `json:"name"`
	Description string                         `json:"description"`
	InputSchema anthropic.ToolInputSchemaParam `json:"input_schema"`
	Function    func(input json.RawMessage) (string, error)
}

// Agent struct represents the core of the AI agent.
type Agent struct {
	client         anthropic.Client
	getUserMessage func() (string, bool)
	tools          []ToolDefinition
	interactive    bool
	systemPrompt   string
}

// NewAgent creates a new Agent instance.
func NewAgent(
	client anthropic.Client,
	getUserMessage func() (string, bool),
	tools []ToolDefinition,
	interactive bool,
	systemPrompt string,
) *Agent {
	return &Agent{
		client:         client,
		getUserMessage: getUserMessage,
		tools:          tools,
		interactive:    interactive,
		systemPrompt:   systemPrompt,
	}
}

// NewClientWithOptions creates a new Anthropic client with the given options.
func NewClientWithOptions(options ...option.RequestOption) anthropic.Client {
	return anthropic.NewClient(options...)
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
		System:    []anthropic.TextBlockParam{{Text: a.systemPrompt}},
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

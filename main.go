package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"tiny-trae/internal/agent"
	"tiny-trae/internal/prompt"
	"tiny-trae/internal/tools"

	"github.com/anthropics/anthropic-sdk-go/option"
)

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
	client := agent.NewClientWithOptions(options...)

	promptFlag := flag.String("p", "", "Accept a string as user input")
	flag.Parse()

	var getUserMessage func() (string, bool)
	var initialMessage string

	if *promptFlag != "" {
		initialMessage = *promptFlag
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

	allTools := tools.GetAllTools()
	systemPrompt := prompt.GetSystemPrompt()
	agentInstance := agent.NewAgent(client, getUserMessage, allTools, *promptFlag == "", systemPrompt)
	err := agentInstance.Run(context.TODO(), initialMessage)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}

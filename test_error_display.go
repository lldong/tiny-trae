package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"tiny-trae/internal/agent"
	"tiny-trae/internal/frontend"
)

func main() {
	// Create a TUI frontend for testing
	tui := frontend.NewTUIFrontend(true)
	
	// Test with a long error message
	longError := strings.Repeat("This is a very long error message that should wrap properly in the TUI display. ", 10)
	
	// Send the error message
	tui.SendMessage(agent.Message{
		Type:    agent.MessageTypeError,
		Content: longError,
	})
	
	// Keep the program running for a few seconds to see the output
	fmt.Println("Testing long error message display...")
	fmt.Println("Error message should wrap properly within the TUI viewport.")
	fmt.Println("Press Ctrl+C to exit.")
	
	// Wait for 10 seconds
	time.Sleep(10 * time.Second)
	
	tui.Close()
	os.Exit(0)
}

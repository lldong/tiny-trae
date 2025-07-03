package profile

import (
	"fmt"
	"strings"

	"tiny-trae/internal/agent"
	"tiny-trae/internal/prompt"
	"tiny-trae/internal/tools"

	"github.com/anthropics/anthropic-sdk-go"
)

// DefaultProfile returns the default profile configuration.
func DefaultProfile() *agent.Profile {
	return &agent.Profile{
		Name:         "default",
		Model:        anthropic.ModelClaudeSonnet4_0,
		MaxTokens:    1024,
		Tools:        tools.GetAllTools(),
		SystemPrompt: prompt.GetSystemPrompt(),
	}
}

// MinimalProfile returns a profile with minimal tools for basic tasks.
func MinimalProfile() *agent.Profile {
	return &agent.Profile{
		Name:         "minimal",
		Model:        anthropic.ModelClaudeSonnet4_0,
		MaxTokens:    1024,
		Tools:        tools.GetMinimalTools(),
		SystemPrompt: prompt.GetMinimalSystemPrompt(),
	}
}

// NewProfile creates a custom profile with the specified configuration.
func NewProfile(name string, model anthropic.Model, maxTokens int64, tools []agent.ToolDefinition, systemPrompt string) *agent.Profile {
	return &agent.Profile{
		Name:         name,
		Model:        model,
		MaxTokens:    maxTokens,
		Tools:        tools,
		SystemPrompt: systemPrompt,
	}
}

// GetAvailableProfiles returns a map of all available built-in profiles.
func GetAvailableProfiles() map[string]*agent.Profile {
	return map[string]*agent.Profile{
		"default": DefaultProfile(),
		"minimal": MinimalProfile(),
	}
}

// ListProfiles prints all available profiles with their descriptions.
func ListProfiles() {
	profiles := GetAvailableProfiles()
	fmt.Println("Available profiles:")
	fmt.Println()

	for name, profile := range profiles {
		var description string
		switch name {
		case "default":
			description = "General-purpose profile with all tools and standard prompt"
		case "minimal":
			description = "Lightweight profile with minimal tools for basic tasks"
		}

		fmt.Printf("  %s:\n", name)
		fmt.Printf("    Description: %s\n", description)
		fmt.Printf("    Model: %s\n", profile.Model)
		fmt.Printf("    Max Tokens: %d\n", profile.MaxTokens)
		fmt.Printf("    Tools: %d available\n", len(profile.Tools))
		fmt.Printf("    System Prompt: %s...\n", strings.TrimSpace(profile.SystemPrompt)[:min(80, len(strings.TrimSpace(profile.SystemPrompt)))])
		fmt.Println()
	}
}

// GetProfileByName returns a profile by its name, or nil if not found.
func GetProfileByName(name string) *agent.Profile {
	profiles := GetAvailableProfiles()
	return profiles[name]
}

// min returns the minimum of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

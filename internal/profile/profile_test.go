package profile

import (
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
)

func TestDefaultProfile(t *testing.T) {
	profile := DefaultProfile()

	if profile.Name != "default" {
		t.Errorf("Expected profile name 'default', got '%s'", profile.Name)
	}

	if profile.Model != anthropic.ModelClaudeSonnet4_0 {
		t.Errorf("Expected model Claude Sonnet 4.0, got %s", profile.Model)
	}

	if profile.MaxTokens != 1024 {
		t.Errorf("Expected max tokens 1024, got %d", profile.MaxTokens)
	}

	if len(profile.Tools) == 0 {
		t.Error("Expected tools to be populated, got empty slice")
	}

	if profile.SystemPrompt == "" {
		t.Error("Expected system prompt to be populated, got empty string")
	}
}

func TestNewProfile(t *testing.T) {
	tools := MinimalProfile().Tools
	profile := NewProfile(
		"test",
		anthropic.ModelClaudeSonnet4_0,
		512,
		tools,
		"Test prompt",
	)

	if profile.Name != "test" {
		t.Errorf("Expected profile name 'test', got '%s'", profile.Name)
	}

	if profile.Model != anthropic.ModelClaudeSonnet4_0 {
		t.Errorf("Expected model Claude Sonnet 4.0, got %s", profile.Model)
	}

	if profile.MaxTokens != 512 {
		t.Errorf("Expected max tokens 512, got %d", profile.MaxTokens)
	}

	if profile.SystemPrompt != "Test prompt" {
		t.Errorf("Expected system prompt 'Test prompt', got '%s'", profile.SystemPrompt)
	}

	if len(profile.Tools) != len(tools) {
		t.Errorf("Expected %d tools, got %d", len(tools), len(profile.Tools))
	}
}

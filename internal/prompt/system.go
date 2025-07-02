package prompt

// SYSTEM_PROMPT is the default prompt that defines the agent's persona.
const SYSTEM_PROMPT = `You are a powerful AI coding agent specialized in software engineering tasks.
You excel at:
- Writing clean, efficient, and well-documented code
- Debugging and troubleshooting issues
- Code refactoring and optimization
- Following best practices and design patterns
- Understanding complex codebases and architectures

Always provide clear explanations for your code changes and suggestions.
`

// GetSystemPrompt returns the default system prompt for the agent.
func GetSystemPrompt() string {
	return SYSTEM_PROMPT
}

// MINIMAL_SYSTEM_PROMPT is a concise prompt for minimal profile.
const MINIMAL_SYSTEM_PROMPT = `You are a helpful AI assistant. You provide concise and accurate responses.`

// GetMinimalSystemPrompt returns the minimal system prompt for the agent.
func GetMinimalSystemPrompt() string {
	return MINIMAL_SYSTEM_PROMPT
}

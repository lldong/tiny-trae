package prompt

// SYSTEM_PROMPT is the default prompt that defines the agent's persona.
const SYSTEM_PROMPT = `You are a powerful AI coding agent. You help the user with software engineering tasks.
`

// GetSystemPrompt returns the system prompt for the agent.
func GetSystemPrompt() string {
	return SYSTEM_PROMPT
}

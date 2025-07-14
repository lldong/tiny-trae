package frontend

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"tiny-trae/internal/agent"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// TUIFrontend implements the Frontend interface for terminal UI interaction using bubbletea
type TUIFrontend struct {
	program     *tea.Program
	model       tuiModel
	inputCh     chan string
	messageCh   chan agent.Message
	interactive bool
	done        chan bool
}

// tuiModel represents the state of the TUI
type tuiModel struct {
	viewport         viewport.Model
	textInput        textinput.Model
	spinner          spinner.Model
	renderer         *glamour.TermRenderer
	messages         []string
	width            int
	height           int
	inputCh          chan string
	messageCh        chan agent.Message
	interactive      bool
	waitingForInput  bool
	waitingForResponse bool
	processingTool   bool
	currentToolName  string
	ready            bool
}

// messageReceivedMsg is sent when a new message is received
type messageReceivedMsg struct {
	msg agent.Message
}

// inputRequestMsg is sent when input is requested
type inputRequestMsg struct{}

// Define styles
var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("magenta")).
		MarginLeft(1)

	userStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("green"))

	assistantStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("cyan"))

	toolStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("yellow"))

	errorStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("196"))

	systemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	inputStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("blue")).
		Padding(0, 1)
)

// NewTUIFrontend creates a new TUI frontend
func NewTUIFrontend(interactive bool) *TUIFrontend {
	inputCh := make(chan string, 1)
	messageCh := make(chan agent.Message, 10)
	done := make(chan bool, 1)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("magenta"))

	textInput := textinput.New()
	textInput.Placeholder = "Type your message here..."
	textInput.CharLimit = 1000
	textInput.Width = 72 // Initial width (80 - 8), will be updated on window resize
	textInput.SetValue("") // Ensure clean initialization

	// Initialize glamour renderer with dark theme (simplified for faster startup)
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		// Fallback to minimal renderer if initialization fails
		renderer, _ = glamour.NewTermRenderer()
	}

	// Initialize viewport with default dimensions
	viewport := viewport.New(80, 20)
	viewport.YPosition = 3

	model := tuiModel{
		viewport:        viewport,
		textInput:       textInput,
		spinner:         s,
		renderer:        renderer,
		inputCh:         inputCh,
		messageCh:       messageCh,
		interactive:     interactive,
		waitingForInput: false,
		waitingForResponse: false,
		processingTool:  false,
		messages:        []string{},
		ready:           true, // Start ready with default dimensions
		width:           80,
		height:          24,
	}

	tui := &TUIFrontend{
		inputCh:     inputCh,
		messageCh:   messageCh,
		interactive: interactive,
		done:        done,
		model:       model,
	}

	if interactive {
		tui.program = tea.NewProgram(model, tea.WithAltScreen())
		go tui.run()
	}

	return tui
}

// run starts the TUI program
func (t *TUIFrontend) run() {
	if _, err := t.program.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
	}
	t.done <- true
}

// Init initializes the TUI model
func (m tuiModel) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		func() tea.Msg {
			// Send a window size message to trigger initialization
			return tea.WindowSizeMsg{Width: 80, Height: 24}
		},
	)
}

// Update handles messages in the TUI
func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		// Update viewport dimensions
		footerHeight := 4
		verticalMarginHeight := footerHeight
		
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - verticalMarginHeight
		
		// Update text input width accounting for border (2) + padding (2)
		// Leave some margin for proper display
		if msg.Width > 8 {
			m.textInput.Width = msg.Width - 8
		}
		
		// Update glamour renderer width only if it's significantly different to avoid unnecessary recreations
		if m.renderer != nil && msg.Width > 20 {
			newRenderer, err := glamour.NewTermRenderer(
				glamour.WithStandardStyle("dark"),
				glamour.WithWordWrap(msg.Width-10), // Leave some margin
			)
			if err == nil {
				m.renderer = newRenderer
			}
		}

	case tea.KeyMsg:
		if !m.interactive {
			switch msg.String() {
			case "ctrl+c":
				os.Exit(0)
			case "q":
				return m, tea.Quit
			}
		}

		if m.waitingForInput && !m.waitingForResponse && !m.processingTool {
			switch msg.String() {
			case "enter":
				input := m.textInput.Value()
				if input != "" {
					m.inputCh <- input
					m.textInput.SetValue("")
					m.textInput.Blur()
					m.waitingForInput = false
					m.waitingForResponse = true
					// Start spinner for response waiting
					cmds = append(cmds, m.spinner.Tick)
				}
				return m, tea.Batch(cmds...)
			case "ctrl+c":
				os.Exit(0)
			}
			m.textInput, cmd = m.textInput.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			switch msg.String() {
			case "ctrl+c":
				os.Exit(0)
			case "q":
				return m, tea.Quit
			}
		}

	case messageReceivedMsg:
		m.addMessage(msg.msg)
		if msg.msg.Type == agent.MessageTypeToolCall {
			m.processingTool = true
			var toolData agent.ToolCallData
			if err := json.Unmarshal(msg.msg.Data, &toolData); err == nil {
				m.currentToolName = toolData.ToolName
			}
			// Start spinner for tool processing
			cmds = append(cmds, m.spinner.Tick)
		} else if msg.msg.Type == agent.MessageTypeToolResult {
			m.processingTool = false
			m.currentToolName = ""
		} else if msg.msg.Type == agent.MessageTypeAssistant {
			// Assistant response received, no longer waiting
			m.waitingForResponse = false
			// Allow free typing again
			m.waitingForInput = true
			m.textInput.Focus()
		}

	case inputRequestMsg:
		m.waitingForInput = true
		m.waitingForResponse = false
		m.textInput.SetValue("") // Clear any residual content
		m.textInput.Focus()

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		// Only continue ticking if we're actively waiting/processing
		if m.waitingForResponse || m.processingTool {
			cmds = append(cmds, cmd)
		}
	}

	// Update viewport
	m.viewport.SetContent(strings.Join(m.messages, "\n"))

	return m, tea.Batch(cmds...)
}

// View renders the TUI
func (m tuiModel) View() string {
	// Footer
	var footer string
	var statusLine string
	
	if m.processingTool {
		statusLine = fmt.Sprintf(" %s Processing tool: %s", m.spinner.View(), m.currentToolName)
	} else if m.waitingForResponse {
		statusLine = fmt.Sprintf(" %s Waiting for response...", m.spinner.View())
	} else if m.interactive {
		statusLine = systemStyle.Render(" Press 'q' or Ctrl+C to quit")
	} else {
		statusLine = systemStyle.Render(" Press 'q' or Ctrl+C to quit")
	}

	// Always show input box, but disable it when waiting for response or processing
	if m.waitingForResponse || m.processingTool {
		// Show disabled input box with muted style
		disabledStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Foreground(lipgloss.Color("240")).
			Padding(0, 1)
		
		// Create a copy of the input to show disabled state
		disabledInput := m.textInput
		disabledInput.Blur()
		inputBox := disabledStyle.Render(disabledInput.View())
		footer = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, inputBox)
	} else {
		// Show normal input box
		inputBox := inputStyle.Render(m.textInput.View())
		footer = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, inputBox)
	}

	// Main view
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.viewport.View(),
		statusLine,
		footer,
	)
}

// wrapText wraps text to fit within the specified width
func wrapText(text string, width int) string {
	if width <= 0 {
		return text
	}
	
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}
	
	var lines []string
	var currentLine strings.Builder
	
	for _, word := range words {
		wordLen := utf8.RuneCountInString(word)
		lineLen := utf8.RuneCountInString(currentLine.String())
		
		// If adding this word would exceed the width, start a new line
		if lineLen+wordLen+1 > width && lineLen > 0 {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
			currentLine.WriteString(word)
		} else {
			if currentLine.Len() > 0 {
				currentLine.WriteString(" ")
			}
			currentLine.WriteString(word)
		}
	}
	
	// Add the last line
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}
	
	return strings.Join(lines, "\n")
}

// addMessage adds a message to the display
func (m *tuiModel) addMessage(msg agent.Message) {
	var formattedMsg string
	timestamp := time.Now().Format("15:04:05")
	
	// Calculate available width for content (account for timestamp, labels, and margins)
	availableWidth := m.width - 12
	if availableWidth < 20 {
		availableWidth = 20
	}

	switch msg.Type {
	case agent.MessageTypeUserInput:
		content := wrapText(msg.Content, availableWidth-6) // Account for prefix
		formattedMsg = fmt.Sprintf("[%s] %s %s", timestamp, userStyle.Render("You:"), content)
	case agent.MessageTypeAssistant:
		// Use glamour to render markdown content from the assistant
		renderedContent, err := m.renderer.Render(msg.Content)
		if err != nil {
			// Fallback to plain text with wrapping if rendering fails
			content := wrapText(msg.Content, availableWidth-6)
			formattedMsg = fmt.Sprintf("[%s] %s %s", timestamp, assistantStyle.Render("Trae:"), content)
		} else {
			// Clean up the rendered content (remove trailing newlines)
			renderedContent = strings.TrimRight(renderedContent, "\n\r")
			// Add timestamp and label
			formattedMsg = fmt.Sprintf("[%s] %s\n%s", timestamp, assistantStyle.Render("Trae:"), renderedContent)
		}
	case agent.MessageTypeToolCall:
		var toolData agent.ToolCallData
		if err := json.Unmarshal(msg.Data, &toolData); err == nil {
			content := wrapText(fmt.Sprintf("Executing %s", toolData.ToolName), availableWidth-6)
			formattedMsg = fmt.Sprintf("[%s] %s %s", timestamp, toolStyle.Render("Tool:"), content)
		} else {
			content := wrapText(msg.Content, availableWidth-6)
			formattedMsg = fmt.Sprintf("[%s] %s %s", timestamp, toolStyle.Render("Tool:"), content)
		}
	case agent.MessageTypeToolResult:
		var toolResult agent.ToolResultData
		if err := json.Unmarshal(msg.Data, &toolResult); err == nil {
			if toolResult.IsError {
				errorText := fmt.Sprintf("%s: %s", toolResult.ToolName, toolResult.Result)
				wrappedError := wrapText(errorText, availableWidth-8)
				formattedMsg = fmt.Sprintf("[%s] %s %s", timestamp, errorStyle.Render("Error"), errorStyle.Render(wrappedError))
			} else {
				// Truncate long results
				result := toolResult.Result
				if len(result) > 200 {
					result = result[:200] + "..."
				}
				content := wrapText(fmt.Sprintf("%s: %s", toolResult.ToolName, result), availableWidth-8)
				formattedMsg = fmt.Sprintf("[%s] %s %s", timestamp, toolStyle.Render("Result"), content)
			}
		} else {
			content := wrapText(msg.Content, availableWidth-6)
			formattedMsg = fmt.Sprintf("[%s] %s %s", timestamp, toolStyle.Render("Result:"), content)
		}
	case agent.MessageTypeError:
		// Wrap error messages to prevent overflow
		wrappedError := wrapText(msg.Content, availableWidth-8)
		formattedMsg = fmt.Sprintf("[%s] %s %s", timestamp, errorStyle.Render("Error:"), errorStyle.Render(wrappedError))
	case agent.MessageTypeSystemInfo:
		content := wrapText(msg.Content, availableWidth-8)
		formattedMsg = fmt.Sprintf("[%s] %s %s", timestamp, systemStyle.Render("System:"), content)
	default:
		content := wrapText(msg.Content, availableWidth-4)
		formattedMsg = fmt.Sprintf("[%s] %s", timestamp, content)
	}

	m.messages = append(m.messages, formattedMsg)
}

// SendMessage sends a message to the TUI for display
func (t *TUIFrontend) SendMessage(msg agent.Message) {
	if t.interactive && t.program != nil {
		t.program.Send(messageReceivedMsg{msg: msg})
	} else {
		// Fallback to stdout for non-interactive mode
		switch msg.Type {
		case agent.MessageTypeAssistant:
			fmt.Printf("Trae: %s\n", msg.Content)
		case agent.MessageTypeError:
			fmt.Printf("Error: %s\n", msg.Content)
		case agent.MessageTypeSystemInfo:
			fmt.Printf("%s\n", msg.Content)
		}
	}
}

// GetUserInput requests user input from the TUI
func (t *TUIFrontend) GetUserInput() (string, bool) {
	if !t.interactive {
		return "", false
	}

	// Send request for input
	if t.program != nil {
		t.program.Send(inputRequestMsg{})
	}

	// Wait for input
	select {
	case input := <-t.inputCh:
		return input, true
	case <-t.done:
		return "", false
	}
}

// IsInteractive returns whether the TUI frontend is in interactive mode
func (t *TUIFrontend) IsInteractive() bool {
	return t.interactive
}

// Close closes the TUI frontend
func (t *TUIFrontend) Close() {
	if t.interactive && t.program != nil {
		t.program.Quit()
		// Don't wait for done channel if we're already done
		select {
		case <-t.done:
			// Program has exited
		default:
			// Program is still running, wait for it to finish
			<-t.done
		}
	}
}

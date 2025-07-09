# Glamour Integration Demo

The TUI frontend now supports **glamour-powered markdown rendering** for assistant messages!

## Features

- **Syntax highlighting** for code blocks
- **Formatted tables** 
- **Styled headings** and text formatting
- **Blockquotes** and lists
- **Links** (displayed with styling)

## Example Usage

When you run the application (TUI is now the default):
```bash
./tiny-trae
```

Assistant messages that contain markdown will be beautifully rendered with:

### Code Blocks
```go
func main() {
    fmt.Println("Hello, World!")
}
```

### Tables
| Feature | Status |
|---------|---------|
| Markdown | ✅ Supported |
| Syntax Highlighting | ✅ Supported |
| Tables | ✅ Supported |

### Lists
- Item 1
- Item 2
- Item 3

> This is a blockquote that will be nicely styled

The glamour renderer automatically adapts to your terminal's color scheme and width!
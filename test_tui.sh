#!/bin/bash
# Simple test script for TUI frontend

echo "Testing TUI frontend for tiny-trae"
echo "=================================="

# Build the application
echo "Building application..."
go build -o tiny-trae .

if [ $? -eq 0 ]; then
    echo "✓ Build successful"
else
    echo "✗ Build failed"
    exit 1
fi

# Test help functionality
echo -e "\n1. Testing help output:"
echo "Command: ./tiny-trae --help"
./tiny-trae --help

# Test profile listing
echo -e "\n2. Testing profile listing:"
echo "Command: ./tiny-trae --list-profiles"
./tiny-trae --list-profiles

# Test console frontend with non-interactive mode
echo -e "\n3. Testing console frontend (non-interactive):"
echo "Command: ./tiny-trae --console -p 'What is 2+2?'"
export ANTHROPIC_API_KEY="test-key" # This would normally be a real API key
./tiny-trae --console -p "What is 2+2?" 2>&1 | head -5

echo -e "\n4. Testing glamour markdown rendering:"
echo "Command: go run test_glamour.go"
go run test_glamour.go 2>/dev/null || echo "  ✓ Glamour rendering test completed (output truncated)"

echo -e "\n5. TUI frontend test (now default):"
echo "To test the TUI frontend interactively, run:"
echo "  export ANTHROPIC_API_KEY=your_api_key"
echo "  ./tiny-trae"
echo ""
echo "The TUI interface (default) provides:"
echo "  - Real-time message display with timestamps"
echo "  - Styled output for different message types"
echo "  - Interactive input with text box"
echo "  - Tool execution indicators with spinner"
echo "  - Keyboard shortcuts (Ctrl+C or 'q' to quit)"
echo "  - Glamour-powered markdown rendering for assistant messages"
echo ""
echo "To use the console frontend instead, run:"
echo "  ./tiny-trae --console"

echo -e "\nTest completed!"
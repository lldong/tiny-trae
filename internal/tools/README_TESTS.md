# Tools Testing Documentation

This directory contains comprehensive unit tests for all tools in the tiny-trae project.

## Test Coverage

Current test coverage: **93.3%** of statements

## Test Files

- `bash_test.go` - Tests for the bash command execution tool
- `edit_file_test.go` - Tests for the file editing tool
- `list_files_test.go` - Tests for the file listing tool
- `read_file_test.go` - Tests for the file reading tool
- `ripgrep_test.go` - Tests for the ripgrep search tool
- `registry_test.go` - Tests for the tool registry

## Running Tests

```bash
# Run all tests
go test ./internal/tools/

# Run tests with verbose output
go test ./internal/tools/ -v

# Run tests with coverage
go test ./internal/tools/ -cover

# Run specific test
go test ./internal/tools/ -run TestBash
```

## Test Features

Each tool test includes:

- **Happy path testing** - Normal operation scenarios
- **Error handling** - Invalid inputs, missing files, etc.
- **Edge cases** - Empty inputs, boundary conditions
- **JSON validation** - Invalid JSON input handling
- **Definition validation** - Tool metadata verification

## Special Considerations

### Ripgrep Tests
- Automatically skips if `rg` command is not available
- Tests various search patterns and options
- Validates output format and content

### File Operation Tests
- Uses temporary directories to avoid affecting the file system
- Cleans up test files automatically
- Tests both relative and absolute paths

### Bash Tests
- Tests command execution and output capture
- Validates error handling for invalid commands
- Cross-platform compatible

## Dependencies

Tests require:
- Go testing framework (built-in)
- `ripgrep` command for ripgrep tests (optional, tests skip if not available)
- Standard Unix commands (`echo`, `pwd`) for bash tests
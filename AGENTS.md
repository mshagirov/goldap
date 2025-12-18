# AGENTS.md

## Build Commands
- `go build` - Build the application
- `go run .` - Run the application directly
- `go fmt ./...` - Format all Go files
- `go vet ./...` - Run static analysis
- `go test ./...` - Run all tests
- `go test -run <TestName>` - Run specific test

## Code Style Guidelines

### Imports
- Group imports into three sections: standard library, third-party, and local packages
- Use absolute imports for local packages (e.g., `github.com/mshagirov/goldap/internal/config`)

### Naming Conventions
- Use PascalCase for exported types, functions, and constants
- Use camelCase for unexported variables and functions
- Use descriptive names (e.g., `LdapApi`, `TableInfo`, `getConfigPath`)

### Error Handling
- Always handle errors explicitly
- Use fmt.Errorf for wrapping errors with context
- Return errors as the last return value

### Types & Structs
- Export struct fields that need to be accessed from other packages
- Use JSON tags for configuration structs
- Define types for domain-specific concepts

### Package Structure
- Keep related functionality in the same package
- Use internal/ for packages not meant for external use
- Main package should be minimal, focusing on application startup
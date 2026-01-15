# AGENTS.md

## Build Commands
- `go build` - Build the application binary
- `go run .` - Run the application directly without building
- `go fmt ./...` - Format all Go files according to Go standards
- `go vet ./...` - Run static analysis to find potential issues
- `go test ./...` - Run all tests in the project
- `go test -run <TestName>` - Run a specific test function
- `go test -v ./...` - Run tests with verbose output
- `go mod tidy` - Clean up module dependencies
- `go mod download` - Download module dependencies

## Development Workflow
- Always run `go fmt ./...` and `go vet ./...` before committing
- Use `go run .` for quick testing during development
- Build with `go build` to create the `goldap` binary
- The binary is ignored by git (see .gitignore)

## Code Style Guidelines

### Import Organization
```go
import (
    // Standard library
    "fmt"
    "os"
    "strconv"
    
    // Third-party packages
    "github.com/charmbracelet/bubbletea"
    "github.com/go-ldap/ldap/v3"
    
    // Local packages (absolute imports)
    "github.com/mshagirov/goldap/internal/config"
    "github.com/mshagirov/goldap/ldapapi"
)
```

### Naming Conventions
- **Exported types/functions/constants**: PascalCase (`LdapApi`, `TableInfo`, `Config`)
- **Unexported variables/functions**: camelCase (`getConfigPath`, `activeTable`)
- **Package names**: lowercase, single word when possible (`config`, `tabs`, `ldapapi`)
- **Constants**: PascalCase for exported, camelCase for unexported
- **Interface names**: often end with `-er` suffix (e.g., `Reader`, `Writer`)

### Error Handling Patterns
- Always handle errors explicitly, never ignore them
- Use `fmt.Errorf` for wrapping errors with context: `fmt.Errorf("DialURL Error: %v", err)`
- Return errors as the last return value: `(result, error)`
- Use descriptive error messages that include the operation name
- For configuration errors, provide helpful guidance to users

### Struct and Type Guidelines
- Export struct fields that need external access using PascalCase
- Use JSON tags for configuration structs: `LdapUrl string \`json:"LDAP_URL"\``
- Define types for domain-specific concepts to improve code clarity
- Use pointer receivers for methods that modify the struct
- Use value receivers for methods that don't modify the struct

### Function and Method Patterns
- Keep functions focused and small
- Use descriptive function names that explain what they do
- For API methods, use clear naming like `ListUsers()`, `GetTableInfo()`
- Use constructor functions for complex initialization: `NewTabsModel()`
- Group related functionality in the same package

### Package Structure
- `main/`: Application entry point, minimal and focused
- `internal/`: Packages not meant for external use
  - `config/`: Configuration management
  - `login/`: Authentication logic
  - `tabs/`: UI tab management
- `ldapapi/`: LDAP API wrapper, exported for external use
- Keep related functionality in the same package
- Avoid circular dependencies between packages

### TUI-Specific Guidelines (Bubbletea)
- Implement `Init()`, `Update()`, and `View()` methods for models
- Use `tea.Model` interface for UI components
- Handle keyboard shortcuts consistently (e.g., `ctrl+c`, `q` for quit)
- Use `tea.Quit` for graceful application termination
- Manage focus state properly for input components

### Constants and Configuration
- Define constants for magic strings and numbers
- Use descriptive constant names: `UserFilter`, `GroupFilter`
- Group related constants together
- Use configuration files for user-customizable settings

### Testing Guidelines
- Place test files in the same package as the code they test
- Use `*_test.go` naming convention
- Write table-driven tests for multiple scenarios
- Test error paths as well as success paths
- Use descriptive test names that explain what is being tested

### Code Documentation
- Exported functions should have godoc comments
- Keep comments concise and focused on the "why" not the "what"
- Use package-level documentation to explain package purpose
- Document complex algorithms or non-obvious code sections
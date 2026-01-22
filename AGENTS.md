# AGENTS.md

## Build Commands
- `go build` - Build the application binary (`goldap`)
- `go run .` - Run the application directly without building
- `go install .` - Install the binary to $GOPATH/bin
- `go fmt ./...` - Format all Go files according to Go standards
- `go vet ./...` - Run static analysis to find potential issues
- `go mod tidy` - Clean up module dependencies
- `go mod download` - Download module dependencies
- `go clean` - Remove object files and cached files
- `go get <package>` - Add dependencies to current module

## Test Commands
- `go test ./...` - Run all tests in the project
- `go test -run <TestName>` - Run a specific test function
- `go test -v ./...` - Run tests with verbose output
- `go test -run TestFunction ./path/to/package` - Run specific test in specific package
- `go test -cover` - Enable code coverage instrumentation
- `go test -race` - Enable data race detection
- `go test -bench=.` - Run benchmarks
- `go test -coverprofile=coverage.out` - Generate coverage profile
- `go test -timeout=30s` - Set test timeout
- `go test -count=1` - Disable test caching

## Advanced Build Options
- `go build -race` - Build with race detector
- `go build -ldflags="-s -w"` - Strip symbols and reduce binary size
- `go build -trimpath` - Remove file system paths from executable
- `go build -x` - Print the commands during build

## Development Workflow
- Always run `go fmt ./...` and `go vet ./...` before committing
- Use `go run .` for quick testing during development
- Build with `go build` to create the `goldap` binary
- The binary is ignored by git (see .gitignore)
- Use `scripts/local-test-server.sh` to start a test LDAP server for development
- Configuration is stored in `~/.goldapconfig.json` (created automatically with example if missing)

## Code Style Guidelines

### Import Organization
- Group imports in three sections with blank lines between: standard library, third-party, local packages
- Use absolute imports for local packages (full module path)
- Sort imports alphabetically within each group
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
- Use `fmt.Errorf` for wrapping errors with context: `fmt.Errorf("DialURL Error; %v", err)`
- Return errors as the last return value: `(result, error)`
- Use descriptive error messages that include the operation name
- For configuration errors, provide helpful guidance to users
- Follow the specific error message format used throughout: `"OperationName Error; %v"`
- LDAP operations always use deferred connection closing: `defer l.Close()`

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

### Testing and UI Development
- Use `scripts/local-test-server.sh` for Docker-based OpenLDAP testing
- Test server creates sample data from `scripts/*.ldif` files
- Mock LDAP connections for unit tests when possible
- UI components use Bubbletea Model interface with `Init()`, `Update()`, `View()` methods
- Handle keyboard shortcuts consistently (ctrl+c, q for quit, etc.)
- Manage focus state properly for input components (login, search, forms)
- Tables are populated via LDAP search with package-defined columns and attributes
- Support dynamic table filtering with search functionality
- Test files should follow `*_test.go` naming convention and be placed in the same package as the code being tested

### Security and Performance
- Never log LDAP credentials; use echo mode masking for password fields
- Use deferred connection closing: `defer l.Close()` for all LDAP operations
- Cache table data between tab switches to avoid redundant LDAP queries
- Validate user input before constructing LDAP filters
- Store config with appropriate permissions (0666 for config files)

### Dependencies and Frameworks
- Bubbletea framework for TUI components with Model interface
- go-ldap/v3 for LDAP protocol implementation
- lipgloss for terminal styling and responsive design
- golang.org/x/term for terminal size detection
- Keep third-party dependencies minimal and well-maintained

### Project-Specific Patterns
- Main loop in `main.go` handles tab/form navigation with `reload_model` flag
- Tab navigation: n/tab for next, p/shift+tab for previous
- Search: / or ? keys; Form mode: enter on table row
- All LDAP operations use same connection pattern with deferred closing
- Configuration uses camelCase JSON keys despite Go struct field naming
- Table names must match switch cases in `GetTableInfo()` function
- Table definitions in `ldapapi/tabfilters.go` must match cases in `GetTableInfo()`
- Column, attribute, and width arrays are synchronized across all table definitions
- Tables use sequential numbering starting from 1
- Configuration file permissions are set to 0666 when created
- Missing configuration triggers automatic example JSON output to guide users

### Development Setup
- Use `scripts/local-test-server.sh` for Docker-based OpenLDAP testing
- Test server creates sample data from `scripts/*.ldif` files (0-ous.ldif, 1-uids.ldif)
- Test config: localhost:389, base DN "dc=goldap,dc=sh"
- Admin credentials: "cn=admin,dc=goldap,dc=sh" with "admin123"
- Application requires Go 1.25.4+ (see go.mod for current requirement)
- Use `go run .` for debugging; binary `goldap` is ignored by git
- Configuration is stored in `~/.goldapconfig.json` (created automatically with example if missing)

### Additional Guidelines
- Function signatures returning errors follow pattern: `(result, error)` or just `error`
- All LDAP operations must use `defer l.Close()` immediately after successful connection
- Error messages use consistent format: `"OperationName Error; %v"` with semicolon separator
- Configuration methods follow `Set*` naming pattern: `SetUrl()`, `SetBaseDn()`, `SetAdminDn()`
- Table operations use sequential numbering starting from 1 (Table1, Table2, etc.)

## Lint and Quality Assurance
- Run `go fmt ./...` before every commit to ensure consistent formatting
- Run `go vet ./...` to catch potential issues and suspicious constructs
- Use `go test -race ./...` to detect data races in concurrent code
- Consider using `golangci-lint` for comprehensive static analysis (if added to project)
- Ensure all error handling follows the established patterns with descriptive messages
- Test files should follow `*_test.go` naming convention and be placed in the same package as the code being tested
- Use table-driven tests for comprehensive test coverage
- Mock external dependencies (like LDAP connections) in unit tests
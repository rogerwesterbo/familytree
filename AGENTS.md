# Agent Guidelines for familytree

## Build/Test Commands
- `go build ./...` - Build all packages
- `go test ./...` - Run all tests
- `go test ./path/to/package -run TestName` - Run a single test
- `go test -v ./...` - Run tests with verbose output
- `go test -cover ./...` - Run tests with coverage

## Code Style
- **Formatting**: Use `gofmt` or `goimports` for all files
- **Imports**: Group stdlib, external, then local imports with blank lines between
- **Naming**: Use camelCase for unexported, PascalCase for exported; avoid underscores
- **Types**: Prefer explicit types; use type inference only when obvious
- **Error Handling**: Always check errors; wrap with context using `fmt.Errorf("context: %w", err)`
- **Comments**: Document all exported functions, types, and packages with godoc format
- **Testing**: Name tests `TestFunctionName`, use table-driven tests for multiple cases
- **Nil Checks**: Always validate pointer/interface parameters before use

## Project Standards
- Use standard Go project layout: `cmd/`, `pkg/`, `internal/` directories
- Keep functions focused and under 50 lines when possible
- Prefer composition over inheritance; use interfaces for abstraction

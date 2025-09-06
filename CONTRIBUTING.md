# Contributing to DocLoom

Thank you for your interest in contributing to DocLoom! We value all contributions, whether they're bug reports, feature requests, documentation improvements, or code contributions.

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct. Please treat all contributors with respect and professionalism.

## How to Contribute

### Reporting Bugs

Before reporting a bug, please:
1. Check existing [issues](https://github.com/karolswdev/docloom/issues) to avoid duplicates
2. Verify the bug still exists in the latest version

When reporting, include:
- Clear, descriptive title
- Steps to reproduce
- Expected behavior
- Actual behavior
- System information (OS, Go version, DocLoom version)
- Relevant logs or error messages

### Suggesting Features

We welcome feature suggestions! Please:
1. Check if the feature has already been requested
2. Provide a clear use case
3. Explain how it benefits DocLoom users
4. Consider implementation complexity

### Contributing Code

#### Getting Started

1. **Fork the repository**
   ```bash
   git clone https://github.com/karolswdev/docloom.git
   cd docloom
   ```

2. **Set up your development environment**
   - Install Go 1.22 or later
   - Install golangci-lint: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
   - Run `make test` to verify your setup

3. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

#### Development Process

1. **Make your changes**
   - Follow our coding standards (see below)
   - Write tests for new functionality
   - Update documentation as needed

2. **Test your changes**
   ```bash
   make test        # Run unit tests
   make lint        # Check code style
   make ci          # Run all checks
   ```

3. **Commit your changes**
   We use [Conventional Commits](https://www.conventionalcommits.org/):
   ```bash
   git commit -m "feat: add new template type"
   git commit -m "fix: resolve PDF parsing issue"
   git commit -m "docs: update installation guide"
   ```

4. **Push and create a Pull Request**
   ```bash
   git push origin feature/your-feature-name
   ```
   Then create a PR on GitHub with:
   - Clear description of changes
   - Link to related issues
   - Test results or screenshots if applicable

### Coding Standards

#### Go Code Style

- Follow standard Go conventions and idioms
- Use `gofmt` and `goimports` for formatting
- Write clear, self-documenting code
- Add comments for complex logic
- Keep functions small and focused

#### Testing

- Write table-driven tests where appropriate
- Aim for good coverage on critical paths
- Include both positive and negative test cases
- Use descriptive test names

Example test structure:
```go
func TestFeatureName(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "test", "expected", false},
        {"invalid input", "", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := YourFunction(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("YourFunction() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("YourFunction() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

#### Error Handling

- Always check and handle errors
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Use custom error types for domain-specific errors
- Never panic in library code

#### Documentation

- Document all exported functions, types, and packages
- Include examples in documentation where helpful
- Keep README and other docs up to date
- Use clear, concise language

### Review Process

1. **Automated checks** - All CI checks must pass
2. **Code review** - At least one maintainer review required
3. **Testing** - Changes must include appropriate tests
4. **Documentation** - Updates must include necessary documentation

### Quick Contribution Checklist

Before submitting your PR, ensure:

- [ ] Code follows Go conventions
- [ ] All tests pass (`make test`)
- [ ] Linting passes (`make lint`)
- [ ] New features have tests
- [ ] Documentation is updated
- [ ] Commit messages follow conventions
- [ ] PR description is clear and complete

## Development Tips

### Running Tests

```bash
# Run all tests
make test

# Run specific package tests
go test ./internal/ai/...

# Run with coverage
make coverage

# Run with race detection
go test -race ./...
```

### Debugging

```bash
# Run with verbose output
docloom generate --verbose ...

# Use dry-run to preview without API calls
docloom generate --dry-run ...

# Enable debug logging
export DOCLOOM_VERBOSE=true
```

### Building

```bash
# Build for development
make build

# Build with version info
make release

# Cross-compile for multiple platforms
make release-all
```

## Getting Help

- **Questions**: Open a [Discussion](https://github.com/karolswdev/docloom/discussions)
- **Bugs**: Open an [Issue](https://github.com/karolswdev/docloom/issues)
- **Security**: Email the maintainers privately

## Recognition

Contributors are recognized in:
- The [Contributors](https://github.com/karolswdev/docloom/graphs/contributors) page
- Release notes for significant contributions
- The README acknowledgments section

Thank you for helping make DocLoom better!
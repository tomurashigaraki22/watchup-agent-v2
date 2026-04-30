# Contributing to WatchUp Agent

Thank you for your interest in contributing to the WatchUp Agent! This document provides guidelines for contributing to this open source project.

## 🌟 How to Contribute

### 1. **Reporting Issues**
- Use GitHub Issues to report bugs or request features
- Provide detailed information including:
  - Operating system and version
  - Go version
  - Agent version
  - Steps to reproduce the issue
  - Expected vs actual behavior

### 2. **Feature Requests**
- Check existing issues to avoid duplicates
- Describe the use case and benefits
- Consider backward compatibility
- Provide implementation suggestions if possible

### 3. **Code Contributions**
- Fork the repository
- Create a feature branch (`git checkout -b feature/amazing-feature`)
- Make your changes
- Add tests for new functionality
- Ensure all tests pass
- Update documentation as needed
- Commit your changes (`git commit -m 'Add amazing feature'`)
- Push to the branch (`git push origin feature/amazing-feature`)
- Open a Pull Request

## 🔧 Development Setup

### Prerequisites
- Go 1.19 or later
- Git

### Local Development
```bash
# Clone the repository
git clone https://github.com/watchup/watchup-agent.git
cd watchup-agent

# Install dependencies
go mod download

# Build the agent
go build -o watchup-agent cmd/agent/main.go cmd/agent/setup.go

# Run tests
go test ./...
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/metrics/
```

## 📋 Code Guidelines

### Code Style
- Follow standard Go conventions
- Use `gofmt` to format code
- Use meaningful variable and function names
- Add comments for exported functions and complex logic

### Commit Messages
- Use clear, descriptive commit messages
- Start with a verb in present tense
- Keep the first line under 50 characters
- Add detailed description if needed

Example:
```
Add CPU temperature monitoring

- Implement temperature collection for Linux systems
- Add configuration option for temperature monitoring
- Include tests for temperature metrics
```

### Pull Request Guidelines
- Keep PRs focused on a single feature or fix
- Include tests for new functionality
- Update documentation as needed
- Ensure CI passes
- Respond to review feedback promptly

## 🧪 Testing Guidelines

### Unit Tests
- Write tests for all new functions
- Use table-driven tests where appropriate
- Mock external dependencies
- Aim for good test coverage

### Integration Tests
- Test complete workflows
- Use temporary files for file operations
- Clean up resources after tests

### Example Test
```go
func TestGetCPUUsage(t *testing.T) {
    metrics, err := GetCPUUsage()
    if err != nil {
        t.Fatalf("GetCPUUsage failed: %v", err)
    }
    
    if metrics.UsagePercent < 0 || metrics.UsagePercent > 100 {
        t.Errorf("Invalid CPU usage: %f", metrics.UsagePercent)
    }
}
```

## 📚 Documentation

### Code Documentation
- Document all exported functions
- Use Go doc conventions
- Include examples for complex functions

### README Updates
- Update README.md for new features
- Keep installation instructions current
- Add configuration examples

### API Documentation
- Update MONITORING_CAPABILITIES.md for new metrics
- Document configuration options
- Include JSON schema examples

## 🔒 Security

### Reporting Security Issues
- **DO NOT** open public issues for security vulnerabilities
- Email security issues to: security@watchup.com
- Include detailed information about the vulnerability
- Allow time for the issue to be addressed before disclosure

### Security Guidelines
- Never commit secrets or credentials
- Use secure coding practices
- Validate all inputs
- Handle errors gracefully

## 🎯 Areas for Contribution

### High Priority
- **New Metrics**: Additional system monitoring capabilities
- **Platform Support**: Improve cross-platform compatibility
- **Performance**: Optimize metric collection performance
- **Documentation**: Improve setup and usage guides

### Medium Priority
- **Configuration**: Enhanced configuration options
- **Logging**: Improved logging and debugging
- **Error Handling**: Better error messages and recovery
- **Testing**: Increase test coverage

### Low Priority
- **UI/UX**: Improve command-line interface
- **Packaging**: Distribution packages (deb, rpm, etc.)
- **Examples**: More configuration examples
- **Integrations**: Third-party integrations

## 🏷️ Versioning

We use [Semantic Versioning](https://semver.org/):
- **MAJOR**: Incompatible API changes
- **MINOR**: New functionality (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

## 📄 License

By contributing to WatchUp Agent, you agree that your contributions will be licensed under the MIT License.

## 🤝 Code of Conduct

### Our Pledge
We are committed to making participation in this project a harassment-free experience for everyone, regardless of age, body size, disability, ethnicity, gender identity and expression, level of experience, nationality, personal appearance, race, religion, or sexual identity and orientation.

### Our Standards
- Use welcoming and inclusive language
- Be respectful of differing viewpoints and experiences
- Gracefully accept constructive criticism
- Focus on what is best for the community
- Show empathy towards other community members

### Enforcement
Instances of abusive, harassing, or otherwise unacceptable behavior may be reported by contacting the project team at conduct@watchup.com.

## 🙋 Getting Help

- **Documentation**: Check README.md and docs/
- **Issues**: Search existing GitHub issues
- **Discussions**: Use GitHub Discussions for questions
- **Community**: Join our community channels (links in README)

## 🎉 Recognition

Contributors will be recognized in:
- CONTRIBUTORS.md file
- Release notes for significant contributions
- Project documentation

Thank you for contributing to WatchUp Agent! 🚀
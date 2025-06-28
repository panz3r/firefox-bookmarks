# Contributing to Firefox Bookmarks Converter

Thank you for your interest in contributing to Firefox Bookmarks Converter! Here are some guidelines to help you get started.

## Development Setup

1. Fork the repository
2. Clone your fork: `git clone git@github.com:yourusername/firefox-bookmarks.git`
3. Create a new branch: `git checkout -b my-feature-branch`

## Building and Testing

```bash
# Build the application
make build

# Run tests
make test

# Build for all platforms
make build-all

# Run performance benchmarks
make benchmark
```

## Pull Requests

1. Push your changes to your fork
2. Submit a pull request to the main repository
3. Ensure your PR has a clear description of the changes and their purpose
4. Make sure all tests pass
5. Update documentation if necessary

## Code Style

Please follow the Go standard code style and formatting guidelines:

- Run `gofmt` before submitting code
- Follow the [Effective Go](https://golang.org/doc/effective_go) guidelines

## Reporting Bugs

When reporting bugs, please include:

- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- Environment details (OS, Go version, etc.)

## Feature Requests

Feature requests are welcome! Please provide:

- A clear description of the feature
- Why it would be useful for the project
- Any implementation ideas you might have

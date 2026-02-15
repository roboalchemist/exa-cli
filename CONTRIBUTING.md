# Contributing

Thanks for your interest in contributing to exa-cli!

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/exa-cli.git`
3. Install dependencies: `make deps`
4. Build: `make build`
5. Run smoke tests: `make test`

## Development

```bash
make build          # Build the binary
make test           # Run smoke tests (no API key needed)
make test-integration  # Run integration tests (requires EXA_API_KEY)
make fmt            # Format code
make lint           # Lint with golangci-lint
make check          # Run fmt, lint, test
```

## Submitting Changes

1. Create a feature branch: `git checkout -b my-feature`
2. Make your changes
3. Run `make check` to verify
4. Commit with a clear message
5. Push and open a pull request

## Integration Tests

Integration tests require an `EXA_API_KEY` environment variable. They hit the live Exa API and are skipped automatically when the key is not set.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

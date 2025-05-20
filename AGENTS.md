# Agent Instructions

This repo is a toolkit for Golang projects. It hosts shared packages and tools used across multiple projects. The goal is to provide a consistent and efficient development experience.

## Development Setup

- Ensure Go 1.24 or newer is installed.
- Run `go mod tidy` to keep `go.mod` in sync.
- Format code with `go fmt ./...` before committing.

## Testing

- CI configuration lives under `.github/workflows`.
- All features should have tests.
- Use ginkgo (v2) and gomega for testing, using dot imports for readability.
- Each package should have a corresponding `suite_test.go` file in the same directory.
- Test package names must be suffixed with `_test`.
- Run `go test ./...` locally before opening a PR.

## Fixing Bugs

When addressing a bug, follow a test-driven development approach:

- **Red** – write a test that reproduces the issue and fails.
- **Green** – implement the minimal fix so the new test passes.
- **Refactor** – clean up the solution while keeping all tests green.

## Pull Requests

- Keep PRs focused and include tests for new behavior.
- Follow the commit message guidelines below.

## Commit Message Format

We follow the [Conventional Commits](https://www.conventionalcommits.org/) standard for commit messages:

```
<type>(<optional scope>): <description>

<optional body>

<optional footer>
```

**Types:**

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Changes that don't affect code functionality (formatting, etc.)
- `refactor`: Code changes that neither fix bugs nor add features
- `test`: Adding or correcting tests
- `chore`: Changes to build process, dependencies, etc.

**Examples:**

```
docs: Enhance AI guidance with interaction workflow and testing policies
feat(core): Add response schema validation
fix(parameter-mapper): Correct path parameter handling
```

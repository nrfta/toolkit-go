# Agent Instructions

This repo is a toolkit for golang projects. It hosts shared packages and tools that are used across multiple projects.
The goal is to provide a consistent and efficient development experience.

## Testing

- Find the CI plan in the .github/workflows folder.
- All features should have tests.
- Use ginkgo (v2) and gomega for testing, using dot imports for readability.
- Each package should have a corresponding test suite file in the same directory.
- Test package name should be suffixed with `_test`.

## Fixing Bugs

When addressing a bug, follow a test-driven development approach:

- Red – Write a test that reproduces the issue and fails.
- Green – Implement the minimal fix so the new test passes.
- Refactor – Clean up the solution while keeping all tests green.

## Other

- Run `go mod tidy` to clear dependency issues.

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

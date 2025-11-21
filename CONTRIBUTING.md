# Contributing to Portcall

Thank you for your interest in contributing! This project is a monorepo containing Go and JavaScript/TypeScript (Next.js, Vite React) applications. Please follow these guidelines to help us maintain a high-quality codebase and smooth collaboration.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How to Contribute](#how-to-contribute)
- [Development Setup](#development-setup)
- [Coding Standards](#coding-standards)
- [Pull Requests](#pull-requests)
- [Issue Reporting](#issue-reporting)
- [Commit Messages](#commit-messages)
- [Branching Model](#branching-model)
- [Code Reviews](#code-reviews)

---

## Code of Conduct

Please be respectful and considerate. Harassment or abusive behavior will not be tolerated.

## How to Contribute

1. Fork the repository and clone your fork.
2. Create a new branch for your feature or bugfix (`feature/your-feature`, `fix/your-bug`).
3. Make your changes, following the coding standards below.
4. Add or update tests as needed.
5. Run all tests and linters locally.
6. Submit a pull request (PR) to the `main` branch.

## Development Setup

### Go Projects

- Install [Go](https://golang.org/doc/install) (version 1.20+ recommended).
- Use `go mod tidy` to manage dependencies.
- Run `go test ./...` to test all packages.
- Use `gofmt` or `go fmt ./...` to format code.

### JavaScript/TypeScript Projects

- Install [Node.js](https://nodejs.org/) (LTS recommended) and [npm](https://www.npmjs.com/).
- Run `npm install` in each frontend directory (e.g., `apps/checkout/frontend`).
- Use `npm run lint` and `npm run test` before submitting changes.
- Use Prettier and ESLint for formatting and linting.

## Coding Standards

### Go

- Follow [Effective Go](https://golang.org/doc/effective_go.html).
- Use clear, descriptive names and comments.
- Keep functions small and focused.
- Write table-driven tests where possible.

### JavaScript/TypeScript

- Use ES6+ features and TypeScript where possible.
- Prefer functional components and hooks in React/Next.js.
- Keep components and files small and focused.
- Use consistent formatting (Prettier, ESLint).

## Pull Requests

- Reference related issues in your PR description (e.g., `Closes #123`).
- Provide a clear summary of your changes.
- Ensure your branch is up to date with `main` before submitting.
- Address all review comments and suggestions.

## Issue Reporting

- Search for existing issues before opening a new one.
- Provide clear, descriptive titles and steps to reproduce.
- Include logs, screenshots, or code samples if helpful.

## Commit Messages

- Use clear, concise messages (e.g., `fix: correct typo in README`).
- Use [Conventional Commits](https://www.conventionalcommits.org/) if possible.

## Branching Model

- Use feature branches (`feature/xyz`), bugfix branches (`fix/xyz`), or hotfix branches (`hotfix/xyz`).
- Do not commit directly to `main`.

## Code Reviews

- All PRs require at least one approval before merging.
- Be constructive and respectful in reviews.
- Suggest improvements, but also acknowledge good work.

---

Thank you for helping make Portcall better!

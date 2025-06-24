# Emoji Remover

A simple, fast, and powerful command-line tool to find and remove emojis from your codebases.

[![Go Report Card](https://goreportcard.com/badge/github.com/user/emoji-remover)](https://goreportcard.com/report/github.com/user/emoji-remover)
[![Go CI](https://github.com/user/emoji-remover/actions/workflows/go.yml/badge.svg)](https://github.com/user/emoji-remover/actions/workflows/go.yml)

The problem it's trying to solve is that AI-generated code from tools like Cursor often has emojis in it, which makes it obvious that it was AI-generated and can look unprofessional. This tool helps you maintain a clean, professional, and emoji-free codebase.

## Installation

### With Go

If you have Go installed, you can install `emoji-remover` with a single command:

```sh
go install github.com/user/emoji-remover@latest
```

### From GitHub Releases (for macOS, Linux, Windows)

You can also download a pre-compiled binary for your operating system from the [latest GitHub release](https://github.com/user/emoji-remover/releases/latest).

### With Homebrew (Coming Soon)

```sh
# brew install user/tap/emoji-remover
```

## Usage

You can run `emoji-remover` on one or more files or directories.

```sh
emoji-remover <path1> [<path2> ...]
```

### Examples

- Remove emojis from a single file:
  ```sh
  emoji-remover path/to/your/file.go
  ```
- Remove emojis from an entire directory recursively:
  ```sh
  emoji-remover .
  ```

### Flags

- `--dry-run`: Show which files contain emojis without actually modifying them.
  ```sh
  emoji-remover --dry-run .
  ```
- `--check`: Exit with a non-zero status code if any emojis are found. This is perfect for CI/CD pipelines or pre-commit hooks.
  ```sh
  emoji-remover --check .
  ```

## Pre-commit Hook

You can use `emoji-remover` as a pre-commit hook to automatically clean your files before you commit them.

1.  First, install the [pre-commit](https://pre-commit.com/) framework:
    ```sh
    pip install pre-commit
    ```
2.  Create a `.pre-commit-config.yaml` file in the root of your repository with the following content. This configuration uses the `--check` flag to ensure no emojis make it into your git history.

    ```yaml
    repos:
      - repo: local
        hooks:
          - id: emoji-remover
            name: Check for emojis
            entry: emoji-remover --check
            language: golang
            files: \.(go|py|js|ts|md|txt)$
            # You might need to install the tool first
            # pre-commit install-hooks
    ```

3.  Install the hook:
    ```sh
    pre-commit install
    ```

Now, `emoji-remover` will run on every commit to ensure your codebase stays emoji-free.

## Development

To build the project from source, clone the repository and run:

```sh
go build -o emoji-remover .
```

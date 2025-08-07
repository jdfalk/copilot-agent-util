<!-- file: README.md -->
<!-- version: 1.0.0 -->
<!-- guid: 56b982b3-90bf-4cc0-bd95-168e7d67ee23 -->

# Copilot Agent Utility

A centralized command execution utility designed to solve VS Code task execution issues and provide consistent logging for Copilot/AI agent operations.

## Overview

This Go application serves as a reliable intermediary between VS Code tasks and system commands, ensuring proper working directory handling, comprehensive logging, and consistent output formatting.

## Features

- **Reliable Command Execution**: Handles working directory changes and environment setup
- **Comprehensive Logging**: Outputs to both terminal and log files with timestamps
- **VS Code Integration**: Designed to work seamlessly with VS Code tasks
- **Cross-Platform**: Works on macOS, Linux, and Windows
- **Extensible**: Easy to add new command types and utilities

## Installation

```bash
go install github.com/jdfalk/copilot-agent-util/cmd/copilot-agent-util@latest
```

Or build from source:

```bash
git clone https://github.com/jdfalk/copilot-agent-util.git
cd copilot-agent-util
go build -o bin/copilot-agent-util cmd/copilot-agent-util/main.go
```

## Usage

```bash
# Basic command execution
copilot-agent-util exec "ls -la"

# Git operations
copilot-agent-util git add .
copilot-agent-util git commit -m "feat: add new feature"
copilot-agent-util git push

# Protocol buffer operations
copilot-agent-util buf generate
copilot-agent-util buf generate --module auth

# File operations
copilot-agent-util file cat README.md
copilot-agent-util file ls src/

# Development tools
copilot-agent-util python run script.py
copilot-agent-util npm install
copilot-agent-util uv run main.py
```

## Command Categories

### File Operations
- `file ls <path>` - List directory contents
- `file cat <file>` - Display file contents
- `file cp <src> <dst>` - Copy files/directories
- `file mv <src> <dst>` - Move/rename files/directories

### Git Operations
- `git add <pattern>` - Add files to staging
- `git commit -m <message>` - Commit changes
- `git push` - Push to remote
- `git push --force-with-lease` - Safe force push
- `git status` - Show repository status

### Protocol Buffers
- `buf generate` - Generate all protocol buffers
- `buf generate --module <name>` - Generate specific module
- `buf lint` - Lint protocol buffer files

### Development Tools
- `python run <script>` - Run Python scripts
- `python build` - Build Python projects
- `uv run <command>` - Run commands with uv
- `npm install` - Install npm dependencies
- `npm run <script>` - Run npm scripts
- `npx <command>` - Execute npm packages

## Configuration

The utility reads configuration from:
- `~/.config/copilot-agent-util/config.yaml`
- Environment variables
- Command-line flags

## Logging

All operations are logged to:
- Terminal (stdout/stderr)
- Log files in `./logs/` directory
- Structured JSON logs for automation

## VS Code Integration

Update your `.vscode/tasks.json`:

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Buf Generate",
      "type": "shell",
      "command": "copilot-agent-util buf generate",
      "group": "build",
      "options": {
        "cwd": "${workspaceFolder}"
      }
    }
  ]
}
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

See [LICENSE](LICENSE) file.

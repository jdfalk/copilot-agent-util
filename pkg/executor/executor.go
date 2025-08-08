// file: pkg/executor/executor.go
// version: 1.0.0
// guid: 45f9f040-7810-4983-bd7b-3fc2ae6af814

package executor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// ExecutorConfig holds configuration for command execution
type ExecutorConfig struct {
	WorkingDir string
	LogDir     string
	Verbose    bool
}

// DefaultConfig returns default configuration
func DefaultConfig() *ExecutorConfig {
	wd, _ := os.Getwd()
	return &ExecutorConfig{
		WorkingDir: wd,
		LogDir:     filepath.Join(wd, "logs"),
		Verbose:    false,
	}
}

// EnsureLogDir creates the log directory if it doesn't exist
func (e *ExecutorConfig) EnsureLogDir() error {
	if _, err := os.Stat(e.LogDir); os.IsNotExist(err) {
		return os.MkdirAll(e.LogDir, 0755)
	}
	return nil
}

// NewExecCommand creates the exec command
func NewExecCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "exec [command]",
		Short: "Execute arbitrary commands",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			executeCommand(args[0], args[1:]...)
		},
	}
}

// NewGitCommand creates git subcommands
func NewGitCommand() *cobra.Command {
	gitCmd := &cobra.Command{
		Use:   "git",
		Short: "Git operations",
	}

	// git add
	gitCmd.AddCommand(&cobra.Command{
		Use:   "add [files...]",
		Short: "Add files to staging area",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				args = []string{"."}
			}
			gitArgs := append([]string{"add"}, args...)
			executeCommand("git", gitArgs...)
		},
	})

	// git commit
	gitCmd.AddCommand(&cobra.Command{
		Use:   "commit",
		Short: "Commit changes",
		Run: func(cmd *cobra.Command, args []string) {
			message, _ := cmd.Flags().GetString("message")
			if message == "" {
				message = "feat: automated commit via copilot-agent-util"
			}
			executeCommand("git", "commit", "-m", message)
		},
	})

	// git push
	gitCmd.AddCommand(&cobra.Command{
		Use:   "push",
		Short: "Push to remote repository",
		Run: func(cmd *cobra.Command, args []string) {
			forceWithLease, _ := cmd.Flags().GetBool("force-with-lease")
			if forceWithLease {
				executeCommand("git", "push", "--force-with-lease")
			} else {
				executeCommand("git", "push")
			}
		},
	})

	// Add flags
	gitCmd.PersistentFlags().StringP("message", "m", "", "Commit message")
	gitCmd.PersistentFlags().Bool("force-with-lease", false, "Force push with lease")

	return gitCmd
}

// NewBufCommand creates buf subcommands
func NewBufCommand() *cobra.Command {
	bufCmd := &cobra.Command{
		Use:   "buf",
		Short: "Protocol buffer operations",
	}

	bufCmd.AddCommand(&cobra.Command{
		Use:   "generate",
		Short: "Generate protocol buffers",
		Run: func(cmd *cobra.Command, args []string) {
			module, _ := cmd.Flags().GetString("module")
			if module != "" {
				executeCommand("buf", "generate", "--path", fmt.Sprintf("pkg/%s/proto", module))
			} else {
				executeCommand("buf", "generate")
			}
		},
	})

	bufCmd.AddCommand(&cobra.Command{
		Use:   "lint",
		Short: "Lint protocol buffers",
		Run: func(cmd *cobra.Command, args []string) {
			module, _ := cmd.Flags().GetString("module")
			if module != "" {
				executeCommand("buf", "lint", "--path", fmt.Sprintf("pkg/%s/proto", module))
			} else {
				executeCommand("buf", "lint")
			}
		},
	})

	bufCmd.PersistentFlags().String("module", "", "Specific module to generate/lint")

	return bufCmd
}

// NewFileCommand creates file operation subcommands
func NewFileCommand() *cobra.Command {
	fileCmd := &cobra.Command{
		Use:   "file",
		Short: "File operations",
	}

	fileCmd.AddCommand(&cobra.Command{
		Use:   "ls [path]",
		Short: "List directory contents",
		Run: func(cmd *cobra.Command, args []string) {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}
			executeCommand("ls", "-la", path)
		},
	})

	fileCmd.AddCommand(&cobra.Command{
		Use:   "cat [file]",
		Short: "Display file contents",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			executeCommand("cat", args[0])
		},
	})

	return fileCmd
}

// NewPythonCommand creates Python development subcommands
func NewPythonCommand() *cobra.Command {
	pythonCmd := &cobra.Command{
		Use:   "python",
		Short: "Python development tools",
	}

	pythonCmd.AddCommand(&cobra.Command{
		Use:   "run [script]",
		Short: "Run Python script",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			executeCommand("python3", args[0])
		},
	})

	return pythonCmd
}

// NewNpmCommand creates npm/node subcommands
func NewNpmCommand() *cobra.Command {
	npmCmd := &cobra.Command{
		Use:   "npm",
		Short: "npm/node operations",
	}

	npmCmd.AddCommand(&cobra.Command{
		Use:   "install",
		Short: "Install npm dependencies",
		Run: func(cmd *cobra.Command, args []string) {
			executeCommand("npm", "install")
		},
	})

	return npmCmd
}

// executeCommand executes a command with logging and output to both terminal and log file
func executeCommand(name string, args ...string) {
	config := DefaultConfig()
	if err := config.EnsureLogDir(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create log directory: %v\n", err)
		os.Exit(1)
	}

	logFile := filepath.Join(config.LogDir, fmt.Sprintf("%s_%d.log", name, time.Now().Unix()))

	fmt.Printf("Executing: %s %v\n", name, args)
	fmt.Printf("Log file: %s\n", logFile)

	// Create log file
	logFileHandle, err := os.Create(logFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create log file: %v\n", err)
		os.Exit(1)
	}
	defer logFileHandle.Close()

	// Write command info to log
	fmt.Fprintf(logFileHandle, "Command: %s %v\n", name, args)
	fmt.Fprintf(logFileHandle, "Working Directory: %s\n", config.WorkingDir)
	fmt.Fprintf(logFileHandle, "Timestamp: %s\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(logFileHandle, "--- Output ---\n")

	command := exec.Command(name, args...)
	command.Dir = config.WorkingDir

	// Create writers that tee output to both terminal and log file
	stdoutWriter := &TeeWriter{terminal: os.Stdout, logFile: logFileHandle}
	stderrWriter := &TeeWriter{terminal: os.Stderr, logFile: logFileHandle}

	command.Stdout = stdoutWriter
	command.Stderr = stderrWriter

	if err := command.Run(); err != nil {
		errorMsg := fmt.Sprintf("Command failed: %v\n", err)
		fmt.Fprint(os.Stderr, errorMsg)
		fmt.Fprint(logFileHandle, errorMsg)
		os.Exit(1)
	}

	successMsg := "Command completed successfully\n"
	fmt.Print(successMsg)
	fmt.Fprint(logFileHandle, successMsg)
}

// TeeWriter writes to both terminal and log file
type TeeWriter struct {
	terminal *os.File
	logFile  *os.File
}

func (t *TeeWriter) Write(p []byte) (n int, err error) {
	// Write to terminal
	if _, err := t.terminal.Write(p); err != nil {
		return 0, err
	}
	// Write to log file
	return t.logFile.Write(p)
}

// NewLintersCommand creates linting subcommands
func NewLintersCommand() *cobra.Command {
	lintersCmd := &cobra.Command{
		Use:   "lint",
		Short: "Code linting and formatting tools",
	}

	// Prettier
	lintersCmd.AddCommand(&cobra.Command{
		Use:   "prettier [files...]",
		Short: "Format code with Prettier",
		Run: func(cmd *cobra.Command, args []string) {
			check, _ := cmd.Flags().GetBool("check")
			write, _ := cmd.Flags().GetBool("write")

			prettierArgs := []string{}
			if check {
				prettierArgs = append(prettierArgs, "--check")
			}
			if write {
				prettierArgs = append(prettierArgs, "--write")
			}

			if len(args) == 0 {
				args = []string{"."}
			}
			prettierArgs = append(prettierArgs, args...)
			executeCommand("prettier", prettierArgs...)
		},
	})

	// ESLint
	lintersCmd.AddCommand(&cobra.Command{
		Use:   "eslint [files...]",
		Short: "Lint JavaScript/TypeScript with ESLint",
		Run: func(cmd *cobra.Command, args []string) {
			fix, _ := cmd.Flags().GetBool("fix")

			eslintArgs := []string{}
			if fix {
				eslintArgs = append(eslintArgs, "--fix")
			}

			if len(args) == 0 {
				args = []string{"."}
			}
			eslintArgs = append(eslintArgs, args...)
			executeCommand("eslint", eslintArgs...)
		},
	})

	// Black (Python formatter)
	lintersCmd.AddCommand(&cobra.Command{
		Use:   "black [files...]",
		Short: "Format Python code with Black",
		Run: func(cmd *cobra.Command, args []string) {
			check, _ := cmd.Flags().GetBool("check")

			blackArgs := []string{}
			if check {
				blackArgs = append(blackArgs, "--check")
			}

			if len(args) == 0 {
				args = []string{"."}
			}
			blackArgs = append(blackArgs, args...)
			executeCommand("black", blackArgs...)
		},
	})

	// Flake8 (Python linter)
	lintersCmd.AddCommand(&cobra.Command{
		Use:   "flake8 [files...]",
		Short: "Lint Python code with Flake8",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				args = []string{"."}
			}
			executeCommand("flake8", args...)
		},
	})

	// Go fmt
	lintersCmd.AddCommand(&cobra.Command{
		Use:   "gofmt [files...]",
		Short: "Format Go code with gofmt",
		Run: func(cmd *cobra.Command, args []string) {
			write, _ := cmd.Flags().GetBool("write")

			gofmtArgs := []string{}
			if write {
				gofmtArgs = append(gofmtArgs, "-w")
			}

			if len(args) == 0 {
				args = []string{"."}
			}
			gofmtArgs = append(gofmtArgs, args...)
			executeCommand("gofmt", gofmtArgs...)
		},
	})

	// Go vet
	lintersCmd.AddCommand(&cobra.Command{
		Use:   "govet [packages...]",
		Short: "Examine Go code with go vet",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				args = []string{"./..."}
			}
			govetArgs := append([]string{"vet"}, args...)
			executeCommand("go", govetArgs...)
		},
	})

	// Golint
	lintersCmd.AddCommand(&cobra.Command{
		Use:   "golint [packages...]",
		Short: "Lint Go code with golint",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				args = []string{"./..."}
			}
			executeCommand("golint", args...)
		},
	})

	// Buf lint (alias for buf lint command)
	lintersCmd.AddCommand(&cobra.Command{
		Use:   "buf [files...]",
		Short: "Lint protocol buffers with buf",
		Run: func(cmd *cobra.Command, args []string) {
			module, _ := cmd.Flags().GetString("module")
			if module != "" {
				executeCommand("buf", "lint", "--path", fmt.Sprintf("pkg/%s/proto", module))
			} else {
				executeCommand("buf", "lint")
			}
		},
	})

	// Markdownlint
	lintersCmd.AddCommand(&cobra.Command{
		Use:   "markdownlint [files...]",
		Short: "Lint Markdown files",
		Run: func(cmd *cobra.Command, args []string) {
			fix, _ := cmd.Flags().GetBool("fix")

			markdownArgs := []string{}
			if fix {
				markdownArgs = append(markdownArgs, "--fix")
			}

			if len(args) == 0 {
				args = []string{"**/*.md"}
			}
			markdownArgs = append(markdownArgs, args...)
			executeCommand("markdownlint", markdownArgs...)
		},
	})

	// Shellcheck
	lintersCmd.AddCommand(&cobra.Command{
		Use:   "shellcheck [files...]",
		Short: "Lint shell scripts with shellcheck",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				executeCommand("find", ".", "-name", "*.sh", "-exec", "shellcheck", "{}", "+")
			} else {
				executeCommand("shellcheck", args...)
			}
		},
	})

	// Add common flags
	lintersCmd.PersistentFlags().Bool("check", false, "Check only, don't write changes")
	lintersCmd.PersistentFlags().Bool("write", false, "Write changes to files")
	lintersCmd.PersistentFlags().Bool("fix", false, "Fix issues automatically")
	lintersCmd.PersistentFlags().String("module", "", "Specific module for buf lint")

	return lintersCmd
}

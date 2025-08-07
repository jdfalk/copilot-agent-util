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
			config := DefaultConfig()
			config.EnsureLogDir()

			fmt.Printf("Executing: %v\n", args)

			command := exec.Command(args[0], args[1:]...)
			command.Dir = config.WorkingDir
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr

			if err := command.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Command failed: %v\n", err)
				os.Exit(1)
			}
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
			executeGitCommand("add", args...)
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
			executeGitCommand("commit", "-m", message)
		},
	})

	// git push
	gitCmd.AddCommand(&cobra.Command{
		Use:   "push",
		Short: "Push to remote repository",
		Run: func(cmd *cobra.Command, args []string) {
			forceWithLease, _ := cmd.Flags().GetBool("force-with-lease")
			if forceWithLease {
				executeGitCommand("push", "--force-with-lease")
			} else {
				executeGitCommand("push")
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
				executeBufCommand("generate", "--path", fmt.Sprintf("pkg/%s/proto", module))
			} else {
				executeBufCommand("generate")
			}
		},
	})

	bufCmd.PersistentFlags().String("module", "", "Specific module to generate")

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

// Helper functions
func executeGitCommand(args ...string) {
	executeCommand("git", args...)
}

func executeBufCommand(args ...string) {
	executeCommand("buf", args...)
}

func executeCommand(name string, args ...string) {
	config := DefaultConfig()
	config.EnsureLogDir()

	logFile := filepath.Join(config.LogDir, fmt.Sprintf("%s_%d.log", name, time.Now().Unix()))

	fmt.Printf("Executing: %s %v\n", name, args)
	fmt.Printf("Log file: %s\n", logFile)

	command := exec.Command(name, args...)
	command.Dir = config.WorkingDir
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Command failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Command completed successfully")
}

// file: cmd/copilot-agent-util/main.go
// version: 1.0.0
// guid: b0a53ad3-47c2-4374-8159-ec4d416f655b

package main

import (
	"fmt"
	"os"

	"github.com/jdfalk/copilot-agent-util/pkg/executor"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "copilot-agent-util",
	Short: "Centralized utility for Copilot/AI agent command execution",
	Long: `A reliable command execution utility designed to solve VS Code task
execution issues and provide consistent logging for Copilot/AI agent operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(executor.NewExecCommand())
	rootCmd.AddCommand(executor.NewGitCommand())
	rootCmd.AddCommand(executor.NewBufCommand())
	rootCmd.AddCommand(executor.NewFileCommand())
	rootCmd.AddCommand(executor.NewPythonCommand())
	rootCmd.AddCommand(executor.NewNpmCommand())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

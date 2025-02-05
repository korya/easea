package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cobra.Command{
		Use:   "easea",
		Short: "Easea is a tool for quickly opening output of search tools",
	}

	initCmd := cobra.Command{
		Use:   "init",
		Short: "Initialize the tool",
		Args:  cobra.NoArgs,
		Run:   initCmdFn,
	}
	rootCmd.AddCommand(&initCmd)

	uninitCmd := cobra.Command{
		Use:   "uninit",
		Short: "Uninit the tool",
		Args:  cobra.NoArgs,
		Run:   uninitCmdFn,
	}
	rootCmd.AddCommand(&uninitCmd)

	handleCmd := cobra.Command{
		Use:   "handle <FORMAT>",
		Short: "Handle the output of a command",
		Args:  cobra.ExactArgs(1),
		Run:   handleCmdFn,
	}
	handleCmd.Flags().StringP("command", "c", "", "Command to be run")
	rootCmd.AddCommand(&handleCmd)

	openCmd := cobra.Command{
		Use:   "open <SHIFT>",
		Short: "Open the line at the given shift: 1 -- first line, 2 -- second line, -1 -- last line, etc.",
		Args:  cobra.ExactArgs(1),
		Run:   openCmdFn,
	}
	rootCmd.AddCommand(&openCmd)

	log.SetOutput(os.Stderr)

	if err := rootCmd.Execute(); err != nil {
		dief("%s", err)
	}
}

func dief(format string, args ...interface{}) {
	log.Fatalf("Error: "+format, args...)
}

func warnf(format string, args ...interface{}) {
	log.Printf("[W] easea: "+format, args...)
}

func infof(format string, args ...interface{}) {
	log.Printf("[I] easea: "+format, args...)
}

func initCmdFn(cmd *cobra.Command, args []string) {
	infof("Initializing")

	s := loadState()
	if s.IsDefined() {
		if err := s.Validate(); err != nil {
			dief("invalid state file: %s", err)
		}

		if s.IsInitialized() {
			infof("Skipping: already initialized")
			return
		}
	}

	f, err := os.CreateTemp("", "easea-*")
	if err != nil {
		dief("Failed to create temp file: %s", err)
	}
	defer f.Close()

	infof("Created temp file %q", f.Name())

	fmt.Printf(`# Easea configuration
export %s=%q
trap 'easea uninit' EXIT
`, stateFileEnvVar, f.Name())
}

func uninitCmdFn(cmd *cobra.Command, args []string) {
	infof("Uninitializing")

	s := loadState()
	if !s.IsInitialized() {
		infof("Skipping: not initialized")
		return
	}

	if err := os.Remove(s.Filename); err != nil {
		warnf("Failed to remove state file %q: %s", s.Filename, err)
	}

	infof("Removed state file %q", s.Filename)
}

func handleCmdFn(cmd *cobra.Command, args []string) {
	format := args[0]
	command, err := cmd.Flags().GetString("command")
	if err != nil {
		dief("Failed to get command: %s", err)
	}

	log.Printf("Handling format %q with command %q", format, command)
}

func openCmdFn(cmd *cobra.Command, args []string) {
	format := args[0]
	command, err := cmd.Flags().GetString("command")
	if err != nil {
		dief("Failed to get command: %s", err)
	}

	log.Printf("Handling format %q with command %q", format, command)
}

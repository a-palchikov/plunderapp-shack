package cmd

import (
	"fmt"
	"os"

	"github.com/plunder-app/shack/pkg/network"

	"github.com/spf13/cobra"
)

// Release - this struct contains the release information populated when building shack
var Release struct {
	Version string
	Build   string
}

func init() {
	shackCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "shack.yaml", "The path to the shack environment configuration")
	// Main function commands
	shackCmd.AddCommand(shackExample)
	shackCmd.AddCommand(shackNetwork)
	shackCmd.AddCommand(shackVM)
	shackCmd.AddCommand(shackVersion)
}

// shackCmd is the parent command
var shackCmd = &cobra.Command{
	Use:   "shack",
	Short: "This is a tool for building a deployment environment",
}

// Execute - starts the command parsing process
func Execute() {
	if err := shackCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

// // Sub commands
var shackVersion = &cobra.Command{
	Use:   "version",
	Short: "Version and Release information about the shack environment manager",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("shack Release Information")
		fmt.Println("Version:  ", Release.Version)
		fmt.Println("Build:    ", Release.Build)
	},
}

var shackNetwork = &cobra.Command{
	Use:   "network",
	Short: "Manage networking",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var shackExample = &cobra.Command{
	Use:   "example",
	Short: "Print example configuratiopn",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(network.ExampleConfig())
	},
}

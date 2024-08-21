package cmd

import (
	"fmt"
	"os"

	"lumino/core"
	"lumino/logger"

	"github.com/spf13/cobra"
)

// log is the package-level logger instance
var log = logger.NewLogger()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "luminocli",
	Short:   "Lumino CLI is a command line interface for interacting with the Lumino network",
	Version: core.Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// init initializes the root command by setting up persistent flags
func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.luminocli.yaml)")
	rootCmd.PersistentFlags().StringP("rpc-url", "r", "", "RPC URL of the Ethereum node")
}

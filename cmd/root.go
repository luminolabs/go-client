package cmd

import (
	"fmt"
	"os"

	"lumino/core"
	"lumino/logger"

	"github.com/spf13/cobra"
)

var log = logger.NewLogger()

var rootCmd = &cobra.Command{
	Use:     "luminocli",
	Short:   "Lumino CLI is a command line interface for interacting with the Lumino network",
	Version: core.Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.luminocli.yaml)")
	rootCmd.PersistentFlags().StringP("rpc-url", "r", "", "RPC URL of the Ethereum node")
}

func main() {
	Execute()
}

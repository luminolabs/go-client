package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
}

func main() {
	var rootCmd = &cobra.Command{Use: "luminocli"}

	var stakeCmd = &cobra.Command{
		Use:   "stake [amount]",
		Short: "Stake tokens in the Lumino network",
		Args:  cobra.ExactArgs(1),
		Run:   runStake,
	}

	var unstakeCmd = &cobra.Command{
		Use:   "unstake [amount]",
		Short: "Unstake tokens from the Lumino network",
		Args:  cobra.ExactArgs(1),
		Run:   runUnstake,
	}

	rootCmd.AddCommand(stakeCmd, unstakeCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runStake(cmd *cobra.Command, args []string) {
	// Implementation for stake command
	log.Info("Staking tokens")
}

func runUnstake(cmd *cobra.Command, args []string) {
	// Implementation for unstake command
	log.Info("Unstaking tokens")
}

package cmd

import (
	"github.com/spf13/cobra"
)

var unstakeCmd = &cobra.Command{
	Use:   "unstake",
	Short: "Manage unstaking operations",
}

var unstakeInitiateCmd = &cobra.Command{
	Use:   "initiate <amount>",
	Short: "Initiate unstaking of tokens",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Implementation for initiating unstaking
	},
}

var unstakeWithdrawCmd = &cobra.Command{
	Use:   "withdraw",
	Short: "Withdraw unstaked tokens",
	Run: func(cmd *cobra.Command, args []string) {
		// Implementation for withdrawing unstaked tokens
	},
}

func init() {
	unstakeCmd.AddCommand(unstakeInitiateCmd, unstakeWithdrawCmd)
	rootCmd.AddCommand(unstakeCmd)
}

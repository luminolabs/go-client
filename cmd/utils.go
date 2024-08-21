// generate file yet to be modified

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// confirmAction asks for user confirmation before proceeding
func confirmAction(prompt string) bool {
	fmt.Printf("%s (y/n): ", prompt)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false
	}
	return response == "y" || response == "Y"
}

// handleError prints the error and exits the program
func handleError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

// addCommonFlags adds flags that are used across multiple commands
func addCommonFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("rpc-url", "r", "", "RPC URL of the Ethereum node")
	cmd.Flags().StringP("private-key", "k", "", "Private key for transaction signing")
}

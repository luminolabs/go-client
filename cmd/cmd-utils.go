// generate file yet to be modified

package cmd

import (
	"fmt"
	"os"

	"lumino/utils"

	"github.com/ethereum/go-ethereum/ethclient"
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

// This function takes client as a parameter and returns the epoch and state
func (*UtilsStruct) GetEpochAndState(client *ethclient.Client) (uint32, int64, error) {
	epoch, err := protoUtils.GetEpoch(client)
	if err != nil {
		return 0, 0, err
	}
	bufferPercent, err := cmdUtils.GetBufferPercent()
	if err != nil {
		return 0, 0, err
	}
	state, err := protoUtils.GetDelayedState(client, bufferPercent)
	if err != nil {
		return 0, 0, err
	}
	log.Debug("Epoch ", epoch)
	log.Debug("State ", utils.GetStateName(state))
	return epoch, state, nil
}

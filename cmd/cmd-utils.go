// generate file yet to be modified

package cmd

import (
	"fmt"
	"os"
	"strconv"

	"lumino/utils"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
)

// confirmAction asks for user confirmation before proceeding with an action.
// It returns true if the user confirms, false otherwise.
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

// GetEpochAndState retrieves the current epoch and state from the Ethereum client.
// It returns the epoch as uint32, state as int64, and an error if retrieval fails.
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
	log.Debug("State ", utils.UtilsInterface.GetStateName(state))
	return epoch, state, nil
}

// GetStatesAllowed returns a string representation of allowed states.
// It takes a slice of int representing allowed states and returns a formatted string.
func GetStatesAllowed(states []int) string {
	var statesAllowed string
	for i := 0; i < len(states); i++ {
		if i == len(states)-1 {
			statesAllowed = statesAllowed + strconv.Itoa(states[i]) + ":" + utils.UtilsInterface.GetStateName(int64(states[i]))
		} else {
			statesAllowed = statesAllowed + strconv.Itoa(states[i]) + ":" + utils.UtilsInterface.GetStateName(int64(states[i])) + ", "
		}
	}
	return statesAllowed
}

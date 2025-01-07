package cmd

import (
	"lumino/logger"
	"lumino/utils"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// networkInfoCmd provides functionality to query and display the current state
// of the Lumino network, including epoch number, network state, and timestamp.
// This command requires no additional parameters and outputs results in tabular format.
var networkInfoCmd = &cobra.Command{
	Use:   "networkInfo",
	Short: "Network Info details",
	Long: `Provides the Network details like state, epoch etc.

Example:
  ./lumino networkInfo`,
	Run: initialiseNetworkInfo,
}

func initialiseNetworkInfo(cmd *cobra.Command, args []string) {
	cmdUtils.ExecuteNetworkInfo(cmd.Flags())
}

// ExecuteNetworkInfo manages the network information retrieval workflow. This function:
// 1. Loads and validates network configuration
// 2. Establishes connection to the specified RPC provider
// 3. Sets up logging with appropriate parameters
// 4. Retrieves and displays current network state
// Returns early with error if any step fails.
func (*UtilsStruct) ExecuteNetworkInfo(flagSet *pflag.FlagSet) {
	config, err := cmdUtils.GetConfigData()
	utils.CheckError("Error in getting config: ", err)
	log.Debugf("ExecuteNetworkInfo: Config: %+v", config)

	client := protoUtils.ConnectToEthClient(config.Provider)
	logger.SetLoggerParameters(client, "")

	log.Debug("ExecuteNetworkInfo: Calling GetNetworkInfo()...")
	err = cmdUtils.GetNetworkInfo(client)
	utils.CheckError("Error in getting Network info : ", err)

}

// GetNetworkInfo retrieves detailed network state from the blockchain including:
// 1. Current epoch number
// 2. Network state value
// 3. Current timestamp
// Returns error if network communication fails or state retrieval fails.
// Results are displayed in a formatted table for better readability.
func (*UtilsStruct) GetNetworkInfo(client *ethclient.Client) error {
	callOpts := protoUtils.GetOptions()
	networkInfo, err := stateManagerUtils.NetworkInfo(client, &callOpts)
	if err != nil {
		log.Errorf("failed to get network info: %v", err)
		return err
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Epoch", "State", "Timestamp"})
	table.Append([]string{
		strconv.Itoa(int(networkInfo.EpochNumber)),
		strconv.Itoa(int(networkInfo.State)),
		networkInfo.Timestamp.String(),
	})
	table.Render()
	return nil
}

// Initializes the network information command in the CLI.
// No additional flags are required as this is a read-only operation
// that uses the default network configuration.
func init() {
	rootCmd.AddCommand(networkInfoCmd)
}

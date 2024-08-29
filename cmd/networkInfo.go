package cmd

import (
	"fmt"
	"lumino/logger"
	"lumino/utils"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

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

func (*UtilsStruct) ExecuteNetworkInfo(flagSet *pflag.FlagSet) {
	config, err := cmdUtils.GetConfigData()
	utils.CheckError("Error in getting config: ", err)
	log.Debugf("ExecuteNetworkInfo: Config: %+v", config)

	client := protoUtils.ConnectToEthClient(config.Provider)
	if client == nil {
		log.Fatal("Failed to connect to Ethereum client")
		return
	}
	logger.SetLoggerParameters(client, "")

	log.Debug("ExecuteNetworkInfo: Calling GetNetworkInfo()...")
	err = cmdUtils.GetNetworkInfo(client)
	utils.CheckError("Error in getting Network info : ", err)

}

func (*UtilsStruct) GetNetworkInfo(client *ethclient.Client) error {
	callOpts := protoUtils.GetOptions()
	networkInfo, err := stateManagerUtils.NetworkInfo(client, &callOpts)
	if err != nil {
		return fmt.Errorf("failed to get network info: %w", err)
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

func init() {
	rootCmd.AddCommand(networkInfoCmd)
}

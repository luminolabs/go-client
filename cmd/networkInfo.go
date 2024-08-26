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
	logger.SetLoggerParameters(client, "")

	provider, err := flagSetUtils.GetRootStringProvider()
	utils.CheckError("Error in getting provider: ", err)
	log.Debug("ExecuteNetworkInfo: Provider: ", provider)

	log.Debug("ExecuteNetworkInfo: Calling GetNetworkInfo() with arguments provider = ", provider)
	err = cmdUtils.GetNetworkInfo(client, provider)
	utils.CheckError("Error in getting Network info : ", err)

}

func (*UtilsStruct) GetNetworkInfo(client *ethclient.Client, provider string) error {
	callOpts := protoUtils.GetOptions()
	networkInfo, err := stateManagerUtils.NetworkInfo(client, &callOpts, provider)
	if err != nil {
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

func init() {
	rootCmd.AddCommand(networkInfoCmd)

	var (
		Provider string
	)

	networkInfoCmd.Flags().StringP(Provider, "stakerId", "", "provider")
}

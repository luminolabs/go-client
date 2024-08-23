package cmd

import (
	"lumino/logger"
	"lumino/utils"
	"os"
	"strconv"
	"time"

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
	log.Debugf("ExecutenetworkInfo: Config: %+v", config)

	client := protoUtils.ConnectToEthClient(config.Provider)
	logger.SetLoggerParameters(client, "")

	log.Debug("ExecutenetworkInfo: Calling GetNetworkInfo()")
	epoch, state, err := cmdUtils.GetEpochAndState(client)
	utils.CheckError("Error in getting staker info: ", err)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Epoch", "State", "Timestamp"})
	table.Append([]string{
		strconv.Itoa(int(epoch)),
		strconv.Itoa(int(state)),
		time.Now().String(),
	})
	table.Render()
}

func init() {
	rootCmd.AddCommand(networkInfoCmd)
}

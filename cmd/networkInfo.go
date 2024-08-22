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

var stakerInfoCmd = &cobra.Command{
	Use:   "stakerInfo",
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
	log.Debugf("ExecuteStakerinfo: Config: %+v", config)

	client := protoUtils.ConnectToEthClient(config.Provider)
	logger.SetLoggerParameters(client, "")

	log.Debug("ExecuteStakerinfo: Calling GetNetworkInfo()")
	err = cmdUtils.GetNetworkInfo(client)
	utils.CheckError("Error in getting staker info: ", err)

}

// This function provides the staker details like age, stake, maturity etc.
func (*UtilsStruct) GetNetworkInfo(client *ethclient.Client, stakerId uint32) error {
	callOpts := protoUtils.GetOptions()
	stakerInfo, err := stakeManagerUtils.StakerInfo(client, &callOpts, stakerId)
	if err != nil {
		return err
	}
	maturity, err := stakeManagerUtils.GetMaturity(client, &callOpts, stakerInfo.Age)
	if err != nil {
		return err
	}
	epoch, err := protoUtils.GetEpoch(client)
	if err != nil {
		return err
	}
	influence, err := protoUtils.GetInfluenceSnapshot(client, stakerId, epoch)
	if err != nil {
		return err
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Staker Id", "Staker Address", "Stake", "Age", "Maturity", "Influence"})
	table.Append([]string{
		strconv.Itoa(int(stakerInfo.Id)),
		stakerInfo.Address.String(),
		stakerInfo.Stake.String(),
		strconv.Itoa(int(stakerInfo.Age)),
		strconv.Itoa(int(maturity)),
		influence.String(),
	})
	table.Render()
	return nil
}

func init() {
	rootCmd.AddCommand(stakerInfoCmd)

	var (
		StakerId uint32
	)

	stakerInfoCmd.Flags().Uint32VarP(&StakerId, "stakerId", "", 0, "staker id")
}

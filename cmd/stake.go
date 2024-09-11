package cmd

import (
	"context"
	"lumino/core"
	"lumino/core/types"
	"lumino/logger"
	"lumino/pkg/bindings"
	"lumino/utils"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// stakeCmd represents the stake command
var stakeCmd = &cobra.Command{
	Use:   "stake",
	Short: "Stake LUMINO tokens",
	Long: `Stake LUMINO tokens in the Lumino network.

This command allows users to stake their LUMINO tokens in the network.
Staking is required to participate in network operations and earn rewards.

Example:
  ./lumino stake --address 0x1234567890123456789012345678901234567890 --amount 1000 --password mySecurePassword`,
	Run: initializeStake,
}

func initializeStake(cmd *cobra.Command, args []string) {
	cmdUtils.ExecuteStake(cmd.Flags())
}

func (*UtilsStruct) ExecuteStake(flagSet *pflag.FlagSet) {
	config, err := cmdUtils.GetConfigData()
	utils.CheckError("Error in getting config: ", err)
	log.Debugf("ExecuteStake: Config: %+v", config)

	client := protoUtils.ConnectToEthClient(config.Provider)

	address, err := flagSetUtils.GetStringAddress(flagSet)
	utils.CheckError("Error in getting address: ", err)
	log.Debug("ExecuteStake: Address: ", address)

	logger.SetLoggerParameters(client, address)
	log.Debug("Checking to assign log file...")
	protoUtils.AssignLogFile(flagSet)

	log.Debug("Getting password...")
	password := protoUtils.AssignPassword(flagSet)

	// Check the LUMINO balance of the staker
	balance, err := protoUtils.FetchBalance(context.Background(), client, common.HexToAddress(address))
	utils.CheckError("Failed to get LUMINO balance:"+address, err)

	log.Debug("Getting amount in wei...")
	valueInWei, err := cmdUtils.AssignAmountInWei(flagSet)
	utils.CheckError("Error in getting amount: ", err)
	log.Debug("ExecuteStake: Amount in wei: ", valueInWei)

	log.Debug("Checking for sufficient balance...")
	protoUtils.CheckAmountAndBalance(valueInWei, balance)

	// TODO: fetch minStake from contracts in Future
	minStakeBigInt := big.NewInt(int64(core.MinimumStake))

	// Ensure the stake amount meets the minimum requirement
	if valueInWei.Cmp(minStakeBigInt) < 0 {
		logger.Fatal("Stake amount", valueInWei.String(), "is below minimum required", minStakeBigInt.String())
	}

	stakerId, err := protoUtils.GetStakerId(client, address)
	utils.CheckError("Error in getting stakerId: ", err)
	log.Debug("ExecuteStake: Staker Id: ", stakerId)

	txnArgs := types.TransactionOptions{
		Client:         client,
		AccountAddress: address,
		Password:       password,
		Amount:         valueInWei,
		ChainId:        core.ChainID,
		Config:         config,
	}

	log.Debug("ExecuteStake: Calling StakeTokens() for amount: ", txnArgs.Amount)
	stakeTxnHash, err := cmdUtils.StakeTokens(txnArgs)
	utils.CheckError("Stake error: ", err)

	err = protoUtils.WaitForBlockCompletion(txnArgs.Client, stakeTxnHash.String())
	utils.CheckError("Error in WaitForBlockCompletion for stake: ", err)
}

// This function allows the user to stake tokens in the lumino network and returns the hash
func (*UtilsStruct) StakeTokens(txnArgs types.TransactionOptions) (common.Hash, error) {
	epoch, err := protoUtils.GetEpoch(txnArgs.Client)
	if err != nil {
		return common.Hash{0x00}, err
	}
	log.Debug("StakeCoins: Epoch: ", epoch)

	txnArgs.ContractAddress = core.StakeManagerAddress
	txnArgs.MethodName = "stake"
	txnArgs.Parameters = []interface{}{epoch, txnArgs.Amount}
	txnArgs.ABI = bindings.StakeManagerABI
	txnOpts := protoUtils.GetTransactionOpts(txnArgs)
	log.Debugf("Executing Stake transaction with epoch = %d, amount = %d", epoch, txnArgs.Amount)
	tx, err := stakeManagerUtils.Stake(txnArgs.Client, txnOpts, epoch, txnArgs.Amount)
	if err != nil {
		return common.Hash{0x00}, err
	}
	log.Info("Txn Hash: ", transactionUtils.Hash(tx).Hex())
	return transactionUtils.Hash(tx), nil
}

func init() {
	rootCmd.AddCommand(stakeCmd)

	var (
		stakeValue    string // Amount of LUMINO tokens to stake
		stakerAddress string // Address of the staker
		password      string // Password for the staker's account
		IsWei         bool
	)

	stakeCmd.Flags().StringVarP(&stakeValue, "value", "v", "0", "Amount of LUMINO tokens to stake")
	stakeCmd.Flags().StringVarP(&stakerAddress, "address", "a", "", "Address of the staker")
	stakeCmd.Flags().StringVarP(&password, "password", "", "", "Password for the staker's account")
	stakeCmd.Flags().BoolVarP(&IsWei, "weiValue", "", false, "value passed in wei")

	stakeAmountErr := stakeCmd.MarkFlagRequired("value")
	utils.CheckError("Value error: ", stakeAmountErr)
	stakerAddrErr := stakeCmd.MarkFlagRequired("address")
	utils.CheckError("Address error: ", stakerAddrErr)
}
